document.addEventListener('DOMContentLoaded', function() {
    const ws = new WebSocket('ws://' + location.host + '/ws');
    const chartDiv = document.getElementById('chart');
    const responseDiv = document.getElementById('response');
    const commandInput = document.getElementById('command');
    const sendBtn = document.getElementById('send-btn');
    const themeToggle = document.getElementById('theme-toggle');

    // Theme toggle functionality
    themeToggle.addEventListener('click', () => {
        document.body.classList.toggle('dark-mode');
        themeToggle.textContent = document.body.classList.contains('dark-mode') ? 'Light' : 'Dark';
        localStorage.setItem('theme', document.body.classList.contains('dark-mode') ? 'dark' : 'light');
    });

    // Load saved theme
    if (localStorage.getItem('theme') === 'dark') {
        document.body.classList.add('dark-mode');
        themeToggle.textContent = 'Light';
    } else {
        themeToggle.textContent = 'Dark';
    }

    sendBtn.onclick = function() {
        const cmd = commandInput.value.trim();
        if (cmd) {
            ws.send(JSON.stringify({ type: 'command', command: cmd }));
            commandInput.value = '';
        }
    };

    commandInput.addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            sendBtn.click();
        }
    });

    ws.onopen = function() {
        // Initial connection - server will send candles automatically
    };

    ws.onmessage = function(event) {
        const msg = JSON.parse(event.data);
        if (msg.type === 'candles') {
            plotCandles(msg.data);
        } else if (msg.type === 'response') {
            responseDiv.innerText += msg.message + '\n';
            responseDiv.scrollTop = responseDiv.scrollHeight;
        }
    };

    ws.onclose = function() {
        console.log('WebSocket closed');
    };

    function plotCandles(candles) {
        d3.select(chartDiv).selectAll("*").remove();

        const margin = { top: 50, right: 40, bottom: 50, left: 40 };
        const width = chartDiv.clientWidth - margin.left - margin.right;
        const height = 500 - margin.top - margin.bottom;

        const data = candles.map(c => ({
            ts: +c.TimeStamp,
            open: +c.Open,
            high: +c.High,
            low: +c.Low,
            close: +c.Close,
            volume: +c.Volume
        })).sort((a, b) => a.ts - b.ts);

        // Debug data
        console.log("Candle data:", data);
        data.forEach(d => {
            if (isNaN(d.low) || isNaN(d.high) || isNaN(d.open) || isNaN(d.close)) {
                console.warn("Invalid data point:", d);
            }
        });

        const svg = d3.select(chartDiv)
            .append("svg")
            .attr("width", width + margin.left + margin.right)
            .attr("height", height + margin.top + margin.bottom)
            .append("g")
            .attr("transform", `translate(${margin.left},${margin.top})`);

        const xDomain = [d3.min(data, d => d.ts), d3.max(data, d => d.ts)];
        const x = d3.scaleLinear()
            .domain(xDomain)
            .range([0, width]);

        const y = d3.scaleLinear()
            .domain([d3.min(data, d => d.low), d3.max(data, d => d.high)])
            .range([height, 0])
            .nice();

        const xAxis = svg.append("g")
            .attr("class", "x-axis")
            .attr("transform", `translate(0,${height})`)
            .call(d3.axisBottom(x));

        svg.append("g")
            .attr("class", "y-axis")
            .call(d3.axisLeft(y));

        svg.append("text")
            .attr("x", width / 2)
            .attr("y", -20)
            .attr("text-anchor", "middle")
            .text("Sim Candles");

        svg.append("text")
            .attr("transform", "rotate(-90)")
            .attr("y", 8 + -margin.left)
            .attr("x", -height / 2)
            .attr("text-anchor", "middle")
            .text("Price");

        const candleWidth = Math.min(10, (width / data.length) * 0.8);

        // Set outline color based on theme
        const candleOutlineColor = document.body.classList.contains('dark-mode') ? '#ffffff' : '#000000';

        const g = svg.append("g")
            .selectAll("g")
            .data(data)
            .join("g")
            .attr("transform", d => `translate(${x(d.ts)},0)`);

        // Wicks
        g.append("line")
            .attr("x1", 0)
            .attr("x2", 0)
            .attr("y1", d => {
                if (isNaN(d.low) || isNaN(d.high)) {
                    console.warn("Invalid low/high for wick:", d);
                }
                return y(d.low);
            })
            .attr("y2", d => y(d.high))
            .attr("stroke", candleOutlineColor) // Use dynamic outline color
            .attr("stroke-width", 2)
            .attr("stroke-opacity", 1)
            .attr("class", "wick");

        // Bodies
        g.append("rect")
            .attr("class", "candle-rect")
            .attr("x", -candleWidth / 2)
            .attr("y", d => y(Math.max(d.open, d.close)))
            .attr("height", d => Math.abs(y(d.open) - y(d.close)) || 1)
            .attr("width", candleWidth)
            .attr("fill", d => d.open > d.close ? "#800606" : "#0fff0f")
            .attr("stroke", candleOutlineColor) // Add outline to bodies
            .attr("stroke-width", 1);

        // Info box
        const infoBox = d3.select(chartDiv)
            .append("div")
            .attr("id", "info-box")
            .style("position", "absolute")
            .style("display", "none")
            .style("background", "rgba(0, 0, 0, 0.8)")
            .style("color", "#ffffff")
            .style("padding", "8px")
            .style("border-radius", "4px")
            .style("pointer-events", "none");

        const formatDate = d3.timeFormat("%B %-d, %Y");
        const formatValue = d3.format(".5f");
        const formatChange = ((f) => (y0, y1) => f((y1 - y0) / y0))(d3.format("+.2%"));

        g.on("mouseover", function(event, d) {
            d3.select(this).select("rect").attr("stroke", "red");
            infoBox.style("display", "block")
                .html(`Time: ${formatDate(new Date(d.ts))}<br>
                       Open: ${formatValue(d.open)}<br>
                       High: ${formatValue(d.high)}<br>
                       Low: ${formatValue(d.low)}<br>
                       Close: ${formatValue(d.close)} (${formatChange(d.open, d.close)})<br>
                       Volume: ${formatValue(d.volume)}`)
                .style("left", `${event.pageX + 10}px`)
                .style("top", `${event.pageY - 10}px`);
        })
        .on("mousemove", function(event) {
            infoBox.style("left", `${event.pageX + 10}px`)
                .style("top", `${event.pageY - 10}px`);
        })
        .on("mouseout", function() {
            d3.select(this).select("rect").attr("stroke", candleOutlineColor); // Reset to dynamic outline color
            infoBox.style("display", "none");
        });

        let xOffset = 0;
        const drag = d3.drag()
            .on("start", function() {
                d3.select(this).style("cursor", "grabbing");
            })
            .on("drag", function(event) {
                xOffset += event.dx;
                const domainWidth = xDomain[1] - xDomain[0];
                const pixelsPerUnit = width / domainWidth;
                const newDomainShift = -xOffset / pixelsPerUnit;
                x.domain([xDomain[0] + newDomainShift, xDomain[1] + newDomainShift]);
                
                g.attr("transform", d => `translate(${x(d.ts)},0)`);
                xAxis.call(d3.axisBottom(x));
            })
            .on("end", function() {
                d3.select(this).style("cursor", "grab");
            });

        svg.call(drag).style("cursor", "grab");
    }
});