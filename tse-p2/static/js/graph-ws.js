document.addEventListener('DOMContentLoaded', function() {
    const ws = new WebSocket('ws://' + location.host + '/ws');
    const chartDiv = document.getElementById('chart');
    const responseDiv = document.getElementById('response');
    const commandInput = document.getElementById('command');
    const sendBtn = document.getElementById('sendBtn');

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
        const margin = { top: 50, right: 30, bottom: 30, left: 40 };
        const width = 800 - margin.left - margin.right;
        const height = 600 - margin.top - margin.bottom;

        // Create SVG container
        const svg = d3.select(chartDiv)
            .append("svg")
            .attr("viewBox", [0, 0, width + margin.left + margin.right, height + margin.top + margin.bottom])
            .style("background", "#fafafa")
            .append("g")
            .attr("transform", `translate(${margin.left},${margin.top})`);

        const data = candles.map(c => ({
            date: c.TimeStamp,
            open: c.Open,
            high: c.High,
            low: c.Low,
            close: c.Close
        })).sort((a, b) => a.date - b.date);

        // Positional encodings
        const x = d3.scaleBand()
            .domain(data.map(d => d.date.toString()))
            .range([0, width])
            .padding(0.8); // Increased padding for less crowded candles

        const y = d3.scaleLinear()
            .domain([d3.min(data, d => d.low), d3.max(data, d => d.high)])
            .rangeRound([height, margin.top])
            .nice();

        // Append axes
        svg.append("g")
            .attr("transform", `translate(0,${height})`)
            .call(d3.axisBottom(x)
                .tickValues(x.domain().filter((d, i) => i % Math.ceil(data.length / 10) === 0))
                .tickFormat(d3.timeFormat("%Y-%m-%d")))
            .attr("stroke", "black")
            .call(g => g.select(".domain").remove());

        svg.append("g")
            .call(d3.axisLeft(y)
                .tickFormat(d3.format(".2f")))
            .attr("stroke", "black")
            .call(g => g.selectAll(".tick line").clone()
                .attr("stroke-opacity", 0.2)
                .attr("x2", width))
            .call(g => g.select(".domain").remove());

        // Title
        svg.append("text")
            .attr("x", width / 2)
            .attr("y", 0)
            .attr("text-anchor", "middle")
            .attr("fill", "black")
            .text("Candlestick Chart");

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
            .attr("transform", d => `translate(${x(d.date.toString())},0)`);

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
High: ${formatValue(d.high)}`);
    }
});