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
            responseDiv.innerText += msg.data;
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

        /*** Create SVG container // andredumas/ techanJS OHLC ***/
        const svg = d3.select(chartDiv)
            .append("svg")
            .append("width", width + margin.left + margin.right)
            .append("height", height + margin.top + margin.bottom)
            .style("background", "#ffffff")
            .append("g")
            .attr("transform", `translate(${margin.left},${margin.top})`);

        const data = candles.map(c => ({
            date: +c.TimeStamp,
            open: +c.Open,
            high: +c.High,
            low: +c.Low,
            close: +c.Close,
            volume: +c.Volume
        })).sort((a, b) => a.date - b.date);


        // Positional encodings
        const x = d3.scaleLinear()
            .domain(data.map(d => d.date))
            .range([0, width])
            .padding(1);

        const y = d3.scaleLinear()
            .domain([d3.min(data, d => d.low), d3.max(data, d => d.high)])
            .rangeRound([height, margin.top])
            .nice();

        // X AXIS //
        svg.append("g")
            .attr("class", "x-axis")
            .attr("transform", `translate(0,${height})`)
            .call(d3.axisBottom(x));

        // Y AXIS //
        svg.append("g")
            .attr("class", "y-axis")
            .call(d3.asixLeft(y));
        
        // Title
        svg.append("text")
            .attr("x", width / 2)
            .attr("y", 0)
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

        // Create group for each candle
        const g = svg.append("g")
            .attr("stroke-linecap", "round")
            .attr("stroke", "black")
            .selectAll("g")
            .data(data)
            .join("g")
            .attr("transform", d => `translate(${x(d.date)},0)`);

        // Wicks (high-low lines)
        g.append("line")
            .attr("y1", d => y(d.low))
            .attr("y2", d => y(d.high))
            .attr("stroke-width", 1); // Thin wicks

        // Bodies (open-close rectangles)
        g.append("rect")
            .attr("x", -x.bandwidth() / 2) // Center the rectangle
            .attr("y", d => y(Math.max(d.open, d.close)))
            .attr("height", d => Math.abs(y(d.open) - y(d.close)) || 1) // Minimum height for flat candles
            .attr("width", x.bandwidth() * 0.8) // Slightly narrower candles
            .attr("fill", d => d.open > d.close ? d3.schemeSet1[0] : d3.schemeSet1[2]);

        // Tooltips
        const formatDate = d3.timeFormat("%B %-d, %Y");
        const formatValue = d3.format(".2f");
        const formatChange = ((f) => (y0, y1) => f((y1 - y0) / y0))(d3.format("+.2%"));

        g.append("title")
            .text(d => `${formatDate(d.date)}
Open: ${formatValue(d.open)}
Close: ${formatValue(d.close)} (${formatChange(d.open, d.close)})
Low: ${formatValue(d.low)}
High: ${formatValue(d.high)}
Volume: ${formatValue(d.volume)}`);
    }
});
