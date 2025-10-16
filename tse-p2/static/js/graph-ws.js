document.addEventListener('DOMContentLoaded', function() {
    const ws = new WebSocket('ws://' + location.host + '/ws');
    const chartDiv = document.getElementById('chart');
    const responseDiv = document.getElementById('response');
    const commandInput = document.getElementById('command');
    const sendBtn = document.getElementById('send-btn');

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
        // Clear previous chart
        d3.select(chartDiv).selectAll("*").remove();

        // Chart dimensions and margins
        const margin = { top: 50, right: 40, bottom: 50, left: 40 };
        const width = 960 - margin.left - margin.right;
        const height = 500 - margin.top - margin.bottom;

        // Parse data
        const data = candles.map(c => ({
            ts: +c.TimeStamp,
            open: +c.Open,
            high: +c.High,
            low: +c.Low,
            close: +c.Close,
            volume: +c.Volume
        })).sort((a, b) => a.ts - b.ts);

        // Create SVG container
        const svg = d3.select(chartDiv)
            .append("svg")
            .attr("width", width + margin.left + margin.right)
            .attr("height", height + margin.top + margin.bottom)
            .style("background", "#ffffff")
            .append("g")
            .attr("transform", `translate(${margin.left},${margin.top})`);

        // Scales
        const x = d3.scaleLinear()
            .domain([d3.min(data, d => d.ts), d3.max(data, d => d.ts)])
            .range([0, width]);

        const y = d3.scaleLinear()
            .domain([d3.min(data, d => d.low), d3.max(data, d => d.high)])
            .range([height, 0])
            .nice();

        // Axes
        svg.append("g")
            .attr("class", "x-axis")
            .attr("transform", `translate(0,${height})`)
            .call(d3.axisBottom(x));

        svg.append("g")
            .attr("class", "y-axis")
            .call(d3.axisLeft(y));

        // Title
        svg.append("text")
            .attr("x", width / 2)
            .attr("y", -20)
            .attr("text-anchor", "middle")
            .attr("fill", "black")
            .text("Sim Candles");

        // Y-axis label
        svg.append("text")
            .attr("transform", "rotate(-90)")
            .attr("y", -margin.left + 10)
            .attr("x", -height / 2)
            .attr("fill", "black")
            .attr("text-anchor", "middle")
            .text("Price");

        // Calculate candle width (based on number of candles)
        const candleWidth = Math.min(10, width / data.length * 0.8); // Max 10px, 80% of available space per candle

        // Create group for each candle
        const g = svg.append("g")
            .attr("stroke", "black")
            .selectAll("g")
            .data(data)
            .join("g")
            .attr("transform", d => `translate(${x(d.ts)},0)`);

        // Wicks (high-low lines)
        g.append("line")
            .attr("y1", d => y(d.low))
            .attr("y2", d => y(d.high))
            .attr("stroke", "black")
            .attr("stroke-width", 1);

        // Bodies (open-close rectangles)
        g.append("rect")
            .attr("x", -candleWidth / 2)
            .attr("y", d => y(Math.max(d.open, d.close)))
            .attr("height", d => Math.abs(y(d.open) - y(d.close)) || 1)
            .attr("width", candleWidth)
            .attr("fill", d => d.open > d.close ? "#000000" : "#ffffff")
            .attr("stroke", "black")
            .attr("stroke-width", 1);

        // Tooltips
        const formatDate = d3.timeFormat("%B %-d, %Y");
        const formatValue = d3.format(".2f");
        const formatChange = ((f) => (y0, y1) => f((y1 - y0) / y0))(d3.format("+.2%"));

        g.append("title")
        .text(d => `Ts:${d.ts}
Open: ${formatValue(d.open)}
Close: ${formatValue(d.close)} (${formatChange(d.open, d.close)})
Low: ${formatValue(d.low)}
High: ${formatValue(d.high)}
Volume: ${formatValue(d.volume)}`);
    }
});