document.addEventListener('DOMContentLoaded', function () {
    const ws = new WebSocket('ws://' + location.host + '/ws');
    const chart = document.getElementById('chart');

    // --- UI Elements ---
    const searchInput = document.getElementById('searchInput');
    const searchBtn = document.getElementById('searchBtn');
    const searchResults = document.getElementById('searchResults');
    const subList = document.getElementById('subList');

    // --- State ---
    const priceHistory = new Map(); // key: addr_etype → [{Ts, Open, High, Low, Close}]
    const subscriptions = new Set(); // "addr_etype" strings

    // --- SVG Candle Rendering (your original code, slightly refactored) ---
    const getChartContext = function (candles) {
        const maxVisible = 50;
        const chartWidth = chart.clientWidth || 800;
        const chartHeight = chart.clientHeight || 500;
        const padding = 30;

        if (candles.length === 0) return { in_view: [], ctx: null };

        let time_min = Infinity, time_max = -Infinity;
        let price_min = Infinity, price_max = -Infinity;

        for (const c of candles) {
            time_min = Math.min(time_min, c.Ts);
            time_max = Math.max(time_max, c.Ts);
            price_min = Math.min(price_min, c.Low);
            price_max = Math.max(price_max, c.High);
        }

        const range = price_max - price_min || 1;
        price_min -= range * 0.1;
        price_max += range * 0.1;

        const startIdx = Math.max(0, candles.length - maxVisible);
        const in_view = candles.slice(startIdx);

        return {
            in_view,
            ctx: { time_min, time_max, price_min, price_max, chartWidth, chartHeight, padding }
        };
    };

    const addCandle = function (c, ctx, i, total) {
        const { chartWidth, chartHeight, padding, price_min, price_max } = ctx;
        const usableW = chartWidth - 2 * padding;
        const usableH = chartHeight - 2 * padding;

        const candleW = usableW / total * 0.8;
        const gap = usableW / total * 0.2;
        const x = padding + i * (candleW + gap);

        const y = p => chartHeight - padding - ((p - price_min) / (price_max - price_min)) * usableH;

        const yHigh = y(c.High), yLow = y(c.Low);
        const yOpen = y(c.Open), yClose = y(c.Close);
        const bodyTop = Math.min(yOpen, yClose);
        const bodyHeight = Math.max(Math.abs(yClose - yOpen), 1);

        // Wick
        const wick = document.createElementNS('http://www.w3.org/2000/svg', 'line');
        wick.setAttribute('x1', x + candleW/2); wick.setAttribute('y1', yHigh);
        wick.setAttribute('x2', x + candleW/2); wick.setAttribute('y2', yLow);
        wick.setAttribute('stroke', '#333'); wick.setAttribute('stroke-width', 1);
        chart.appendChild(wick);

        // Body
        const body = document.createElementNS('http://www.w3.org/2000/svg', 'rect');
        body.setAttribute('x', x); body.setAttribute('y', bodyTop);
        body.setAttribute('width', candleW); body.setAttribute('height', bodyHeight);
        body.setAttribute('fill', c.Close >= c.Open ? '#26a69a' : '#ef5350');
        body.setAttribute('stroke', '#000'); body.setAttribute('stroke-width', 0.5);
        chart.appendChild(body);
    };

    const clearChart = () => { while (chart.firstChild) chart.removeChild(chart.firstChild); };

    const renderAll = () => {
        clearChart();
        if (priceHistory.size === 0) return;

        // Merge all subscribed series into one view (or pick first — change as needed)
        let allCandles = [];
        for (const candles of priceHistory.values()) {
            if (candles.length > allCandles.length) allCandles = candles;
        }

        const { in_view, ctx } = getChartContext(allCandles);
        if (!ctx) return;

        in_view.forEach((c, i) => addCandle(c, ctx, i, in_view.length));
    };

    // --- WebSocket Handling ---
    ws.onopen = () => console.log('WebSocket connected');

    ws.onmessage = (event) => {
        const msg = JSON.parse(event.data);

        // Single price update from your Ledger.Emit
        if (msg.address && msg.priceA !== undefined) {
            const key = `${msg.address}_${msg.type || 2}`; // EntityExchange = 2
            if (!subscriptions.has(key)) return;

            let series = priceHistory.get(key);
            if (!series) {
                series = [];
                priceHistory.set(key, series);
            }

            const now = Date.now();
            const price = parseFloat(msg.priceA) || parseFloat(msg.priceB) || 0;

            // Build or update current candle (1-minute buckets)
            const minute = Math.floor(now / 60000) * 60000;
            let current = series[series.length - 1];

            if (!current || current.Ts !== minute) {
                current = { Ts: minute, Open: price, High: price, Low: price, Close: price };
                series.push(current);
            } else {
                current.High = Math.max(current.High, price);
                current.Low = Math.min(current.Low, price);
                current.Close = price;
            }

            renderAll();
        }
    };

    ws.onclose = () => console.log('WebSocket closed');

    // --- Search ---
    searchBtn.onclick = () => {
        const query = searchInput.value.trim();
        if (!query) return;

        ws.send(JSON.stringify({
            type: "search",
            data: { name: query }
        }));
    };

    // --- Handle search results ---
    ws.addEventListener('message', (event) => {
        const msg = JSON.parse(event.data);
        if (msg.type !== "search_results") return;

        searchResults.innerHTML = "<strong>Results:</strong> ";
        if (!msg.data || msg.data.length === 0) {
            searchResults.innerHTML += "Nothing found.";
            return;
        }

        msg.data.forEach(src => {
            const btn = document.createElement('button');
            const key = `${src.address}_${src.etype}`;
            const isSubbed = subscriptions.has(key);

            btn.textContent = isSubbed ? "Unsubscribe" : "Subscribe";
            btn.style.marginLeft = "8px";
            btn.onclick = () => {
                if (isSubbed) {
                    ws.send(JSON.stringify({
                        type: "unsubscribe",
                        data: { addr: src.address, etype: src.etype }
                    }));
                    subscriptions.delete(key);
                } else {
                    ws.send(JSON.stringify({
                        type: "subscribe",
                        data: { name: src.name, addr: src.address, etype: src.etype }
                    }));
                    subscriptions.add(key);
                }
                updateSubList();
                renderAll();
            };

            const div = document.createElement('div');
            div.textContent = `${src.name} (${src.address})`;
            div.appendChild(btn);
            searchResults.appendChild(div);
        });
    });

    // --- Update subscription UI ---
    const updateSubList = () => {
        if (subscriptions.size === 0) {
            subList.innerHTML = "None";
            return;
        }
        subList.innerHTML = "";
        for (const key of subscriptions) {
            const [addr, etype] = key.split('_');
            const span = document.createElement('span');
            span.className = "sub-item";
            span.textContent = `Addr:${addr} Type:${etype}`;
            const unsub = document.createElement('button');
            unsub.textContent = "X";
            unsub.onclick = () => {
                ws.send(JSON.stringify({
                    type: "unsubscribe",
                    data: { addr: parseInt(addr), etype: parseInt(etype) }
                }));
                subscriptions.delete(key);
                updateSubList();
            };
            span.appendChild(unsub);
            subList.appendChild(span);
        }
    };

    // Initial UI
    updateSubList();
    window.addEventListener('resize', renderAll);
});