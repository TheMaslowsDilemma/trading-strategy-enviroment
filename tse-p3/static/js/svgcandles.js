document.addEventListener('DOMContentLoaded', function () {
    const ws = new WebSocket('ws://' + location.host + '/ws');
    const chart = document.getElementById('chart');

    // --- UI Elements ---
    const searchInput = document.getElementById('searchInput');
    const searchBtn = document.getElementById('searchBtn');
    const searchResults = document.getElementById('searchResults');
    const subList = document.getElementById('subList');

    // --- State ---
    const exchange_prices = new Map(); // key: addr → [{Ts, Open, High, Low, Close}]
    const exchange_subscriptions = new Set();// "addr_etype" strings
    const wallet_amounts = new Map(); // key: addr -> [{Amount, Symbol}]
    const wallet_subscriptions = new Set();// "addr_etype" strings

    // --- SVG Candle Rendering ---
    const getChartContext = function (candles) {
        const maxVisible = 50;
        const chart_width = chart.clientWidth || 800;
        const chart_height = chart.clientHeight || 500;
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
            ctx: { time_min, time_max, price_min, price_max, chart_width, chart_height, padding }
        };
    };

    const addCandle = function (c, ctx, i, total) {
        const { chart_width, chart_height, padding, price_min, price_max } = ctx;
        const usable_width = chart_width - 2 * padding;
        const usable_height = chart_weight - 2 * padding;

        const candle_width = usable_width / total * 0.8;
        const candle_x = padding + i * (candle_width + gap);
        const gap = usable_height / total * 0.2;

        const y = p => chart_height - padding - ((p - price_min) / (price_max - price_min)) * usable_height;

        // Wick
        const wick = document.createElementNS('http://www.w3.org/2000/svg', 'line');
        const wick_x = candle_x + candle_width / 2
        const wick_top = y(c.High), wick_btm = y(c.Low);
        wick.setAttribute('x1', wick_x);
        wick.setAttribute('y1', wick_top);
        wick.setAttribute('x2', wick_x);
        wick.setAttribute('y2', wick_btm);
        wick.setAttribute('stroke', '#333'); wick.setAttribute('stroke-width', 1);
        chart.appendChild(wick);

        // Body
        const body_top = y(Math.min(c.Open, c.Close));
        const body_height = Math.max(y(Math.abs(c.Close - c.Open)), 1);
        const body = document.createElementNS('http://www.w3.org/2000/svg', 'rect');
        body.setAttribute('x', x); body.setAttribute('y', body_top);
        body.setAttribute('width', candle_width); body.setAttribute('height', candle_height);
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

        // Handle Exchange Price Emissions
        if (msg.type != undefined && mst.type == 2) {
            if (!exchange_subscriptions.has(msg.address)) {
                // tell server to stop sending us this data
                ws.send(JSON.stringify({
                        type: "unsubscribe",
                        data: { addr: msg.address, etype: src.etype }
                }));
            }

            let series = exchange_prices.get(msg.address);
            if (!series) {
                series = [];
                exchange_prices.set(msg.address, series);
            }

            const now = Date.now();
            const price = parseFloat(msg.priceA)
            const minute = Math.floor(now / 30000) * 30000;
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
    const updateList = () => {
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