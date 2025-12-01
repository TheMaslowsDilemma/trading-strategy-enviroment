/// charts.js

const wallet_etype   = 0;
const exchange_etype = 1;

let charting_state = {
	wallets: {
		subscribed: new Set(),
		in_view:    new Map(),    // addr → div element
	},
	exchanges: {
		subscribed: new Set(),
		in_view:    new Map(),    // addr → { div + svg + candles[]
	},
};

let ws    = null;                   // global ws reference so handlers can use it

function main_handler() {
	ws = new WebSocket('wss://' + location.host + '/ws');

	ws.onopen    = on_open_handler;
	ws.onmessage = on_message_handler;
	ws.onclose   = on_close_handler;

	init_ui_elements();
}

///
/// ------------------- UI Setup ------------------- ///
///

function init_ui_elements() {
	const search_inpt = document.getElementById('search_inpt');
	const search_bttn = document.getElementById('search_bttn');

	search_bttn.onclick = () => {
		const query = search_inpt.value.trim();
		if (!query) return;

		ws.send(JSON.stringify({
			type: "search",
			data: { name: query }
		}));

		search_inpt.value = '';
	};

	search_inpt.addEventListener('keypress', e => {
		if (e.key === 'Enter') {
			search_bttn.click();
		}
	});

	document.addEventListener('click', (e) => {
		const container = document.querySelector('.search-container');
		if (!container.contains(e.target)) {
			document.getElementById('search_rslts').innerHTML = '';
		}
	});
}

///
/// ------------------- Websocket Handlers ------------------- ///
///

function on_open_handler() {
	console.log('WebSocket connection opened');
}

function on_close_handler() {
	console.log('WebSocket closed – clearing everything');
	charting_state = {
		wallets:   { subscribed: new Set(), in_view: new Map() },
		exchanges: { subscribed: new Set(), in_view: new Map() }
	};

	// remove all created divs
	document.querySelectorAll('[id^="wallet-"], [id^="exchange-"]').forEach(el => el.remove());
	clear_chart(chart);
	if (window.update_sub_list) window.update_sub_list();
}

function on_message_handler(event) {
	const msg = JSON.parse(event.data);

	if (!msg.type) {
		console.warn('Message without type field', msg);
		return;
	}

	switch (msg.type) {
		case 'initialize':
			console.log('User initialized:', msg.data);
			break;

		case 'wallet':
			handle_wallet_data(msg.data);
			break;

		case 'exchange':
			handle_exchange_data(msg.data);
			break;

		case 'search_results':
			handle_search_results_data(msg.data);
			break;

		default:
			console.log('Unknown message type:', msg.type);
	}
}

///
/// ------------------- Data Handlers ------------------- ///
///

function handle_wallet_data(wallet) {
	const addr = wallet.address;
	if (!addr || !charting_state.wallets.in_view.has(addr)) return;

	const div           = charting_state.wallets.in_view.get(addr);
	const balance_span  = div.querySelector('.balance');
	const symbol_span   = div.querySelector('.symbol');

	if (balance_span) balance_span.textContent = wallet.balance || '0';
	if (symbol_span)  symbol_span.textContent  = wallet.symbol  || '?';
}

const handle_exchange_data = (exchange) => {
	const addr = exchange?.address;
	if (!addr) {
		console.warn('Exchange missing address');
		return;
	}

	const entry = charting_state.exchanges.in_view.get(addr);
	render_candles(exchange.candles, entry.svg);
}

///
/// ------------------- Search Results UI ------------------- ///
///

function handle_search_results_data(search_results) {
	const container = document.getElementById('search_rslts');
	if (!container) return;

	container.innerHTML = '<strong>Results:</strong><br>';

	if (!search_results || search_results.length === 0) {
		container.innerHTML += 'Nothing found...';
		return;
	}

	search_results.forEach(sr => {
		let already_subbed = is_subscribed(sr.entity_type, sr.address);

		const search_result_div = document.createElement('div');
		search_result_div.setAttribute('class', 'search-result-item')

		const search_result_txt  = document.createElement('span');
		search_result_txt.textContent = `${sr.name} (${sr.address}) `;

		const btn  = document.createElement('button');

		btn.textContent = already_subbed ? 'Unsubscribe' : 'Subscribe';

		btn.onclick = () => {
			already_subbed = is_subscribed(sr.entity_type, sr.address);
			btn.textContent = already_subbed ? 'Unsubscribe' : 'Subscribe';
			if (already_subbed) {
				ws.send(JSON.stringify({ type: 'unsubscribe', data: sr }));
				remove_subscription(sr.entity_type, sr.address);

				if (sr.entity_type === wallet_etype) {
					const el = charting_state.wallets.in_view.get(sr.address);
					if (el) el.remove();
					charting_state.wallets.in_view.delete(sr.address);
				} else {
					const entry = charting_state.exchanges.in_view.get(sr.address);
					if (entry && entry.div) entry.div.remove();
					charting_state.exchanges.in_view.delete(sr.address);
				}
			} else {
				ws.send(JSON.stringify({ type: 'subscribe', data: sr }));
				add_subscription(sr.entity_type, sr.address);

				if (sr.entity_type === wallet_etype) {
					create_wallet_display(sr);
				} else {
					create_exchange_display(sr);
				}
			}

			if (window.update_sub_list) window.update_sub_list();
			btn.textContent = already_subbed ? 'Subscribe' : 'Unsubscribe';
		};

		search_result_div.appendChild(search_result_txt);
		search_result_div.appendChild(btn);
		container.appendChild(search_result_div);
	});
}

///
/// ------------------- Subscription Helpers ------------------- ///

function is_subscribed(entity_type, addr) {
	if (entity_type === wallet_etype)   return charting_state.wallets.subscribed.has(addr);
	if (entity_type === exchange_etype) return charting_state.exchanges.subscribed.has(addr);
	return false;
}

function add_subscription(entity_type, addr) {
	if (entity_type === wallet_etype)   charting_state.wallets.subscribed.add(addr);
	if (entity_type === exchange_etype) charting_state.exchanges.subscribed.add(addr);
}

function remove_subscription(entity_type, addr) {
	if (entity_type === wallet_etype)   charting_state.wallets.subscribed.delete(addr);
	if (entity_type === exchange_etype) charting_state.exchanges.subscribed.delete(addr);
}

///
/// ------------------- Create UI for New Subscriptions ------------------- ///
///

function create_wallet_display(sr) {
	const div = document.createElement('div');
	div.id = `wallet-${sr.address}`;
	div.className = 'wallet-display';

	div.innerHTML = `
		<strong>${sr.name}</strong><br>
		Balance: <span class="balance">loading...</span> <span class="symbol">?</span>
	`;

	// Insert into the LEFT sidebar
	const walletBox = document.querySelector('.wallet-box');
	walletBox.appendChild(div);

	charting_state.wallets.in_view.set(sr.address, div);
}

function create_exchange_display(sr) {
	// Split "TOKENA:e:TOKENB" into [TOKENA, TOKENB]
	const symbols = sr.name.split(':e:');
	const symbolA = symbols[0].trim();  // e.g., "USDC"
	const symbolB = symbols[1].trim();  // e.g., "SOL"

	const div = document.createElement('div');
	div.className = 'exchange-display';

	// Display pair name
	const nameElement = document.createElement('strong');
	nameElement.textContent = sr.name;
	div.appendChild(nameElement);
	div.appendChild(document.createElement('br'));

	// Amount input
	const amountInput = document.createElement('input');
	amountInput.type = 'number';
	amountInput.placeholder = 'Amount';
	amountInput.className = 'amount-input';
	amountInput.min = '0';
	amountInput.step = 'any';
	div.appendChild(amountInput);

	// Buy Button: Buy A with B → swap from B to A
	const buyButton = document.createElement('button');
	buyButton.className = 'buy-bttn';
	buyButton.textContent = `Buy ${symbolA}`;
	buyButton.onclick = () => {
		const amount = parseFloat(amountInput.value);
		if (isNaN(amount) || amount <= 0) {
			alert('Please enter a valid amount');
			return;
		}
		ws.send(JSON.stringify({
			type: "swap",
			data: {
				from_symbol: symbolB,
				to_symbol: symbolA,
				amount_in: amount,
			}
		}));
		amountInput.value = ''; // optional: clear input
	};
	div.appendChild(buyButton);

	// Sell Button: Sell A for B → swap from A to B
	const sellButton = document.createElement('button');
	sellButton.className = 'sell-bttn';
	sellButton.textContent = `Sell ${symbolA}`;
	sellButton.onclick = () => {
		const amount = parseFloat(amountInput.value);
		if (isNaN(amount) || amount <= 0) {
			alert('Please enter a valid amount');
			return;
		}
		ws.send(JSON.stringify({
			type: "swap",
			data: {
				from_symbol: symbolA,
				to_symbol: symbolB,
				amount_in: amount,
			}
		}));
		amountInput.value = '';
	};
	div.appendChild(sellButton);

	// SVG container for chart
	const svg = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
	svg.setAttribute('width', '100%');
	svg.setAttribute('height', '200'); // or desired height
	svg.id = `exchange-${sr.address}`;
	div.appendChild(svg);

	// Append to chart box
	const chartBox = document.querySelector('.chart-box');
	if (chartBox) {
		chartBox.appendChild(div);
	} else {
		console.error('.chart-box not found in DOM');
	}

	// Store in charting state for later access (e.g. updating charts)
	if (!charting_state.exchanges.in_view) {
		charting_state.exchanges.in_view = new Map();
	}
	charting_state.exchanges.in_view.set(sr.address, {
		svg: svg,
		div: div,
		symbolA,
		symbolB,
		pairName: sr.name
	});
}

///
/// ------------------- Candlestick Charting ------------------- ///
///

function get_chart_context(candles, svg) {
	const max_visible = 50;
	const width = svg.clientWidth;
	const height = svg.clientHeight;
	const padding     = 40;

	if (candles.length === 0) return { in_view: [], ctx: null };

	let t_min = Infinity, t_max = -Infinity;
	let p_min = Infinity, p_max = -Infinity;

	candles.forEach( (c) => {
		t_min = Math.min(t_min, c.Ts);
		t_max = Math.max(t_max, c.Ts);
		p_min = Math.min(p_min, c.Low);
		p_max = Math.max(p_max, c.High);
	});

	const price_range = p_max - p_min || 1;
	p_min -= price_range * 0.06;
	p_max += price_range * 0.06;

	const start_idx = Math.max(0, candles.length - max_visible);
	const in_view   = candles.slice(start_idx);

	const usable_w  = width  - 2 * padding;
	const usable_h  = height - 2 * padding;
	const bar_w     = usable_w / in_view.length * 0.8;
	const gap       = usable_w / in_view.length * 0.2;

	const y_scale = p => height - padding - ((p - p_min) / (p_max - p_min)) * usable_h;

	return { in_view, ctx: { y_scale, padding, bar_w, gap } };
}

function add_candle(candle, ctx, idx, svg) {
	const { y_scale, padding, bar_w, gap } = ctx;

	const x      = padding + idx * (bar_w + gap);
	const wick_x = x + bar_w / 2;

	const high_y = y_scale(candle.High);
	const low_y  = y_scale(candle.Low);
	const open_y = y_scale(candle.Open);
	const close_y= y_scale(candle.Close);

	const top    = Math.min(open_y, close_y);
	const bottom = Math.max(open_y, close_y);

	// wick
	const wick = document.createElementNS('http://www.w3.org/2000/svg', 'line');
	wick.setAttribute('x1', wick_x);
	wick.setAttribute('y1', high_y);
	wick.setAttribute('x2', wick_x);
	wick.setAttribute('y2', low_y);
	wick.setAttribute('stroke', 'black');
	wick.setAttribute('stroke-width', 1.5);
	svg.appendChild(wick);

	// body
	const body = document.createElementNS('http://www.w3.org/2000/svg', 'rect');
	body.setAttribute('x', x);
	body.setAttribute('y', top);
	body.setAttribute('width', bar_w);
	body.setAttribute('height', Math.max(bottom - top, 2));
	body.setAttribute('fill', candle.Close >= candle.Open ? '#22a69aff' : '#fe5300ff');
	svg.appendChild(body);
}

function clear_chart(svg) {
	while (svg && svg.firstChild) {
		svg.removeChild(svg.firstChild);
	}
}

function render_candles(candles, svg) {
	if (!svg) {
		console.warn("svg dne");
		return;
	}

	const temp = get_chart_context(candles, svg);
	if (!temp.ctx) {
		clear_chart(svg);
		return;
	}

	const { in_view, ctx } = temp;

	clear_chart(svg);
	in_view.forEach((c, i) => add_candle(c, ctx, i, svg));
}

///
/// ------------------- Start ------------------- ///
///

document.addEventListener('DOMContentLoaded', main_handler);