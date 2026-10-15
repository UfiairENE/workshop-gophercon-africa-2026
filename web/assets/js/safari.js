(function () {
	let ctx;
	const buffers = {};
	const loading = {};

	function audio() {
		const AC = window.AudioContext || window.webkitAudioContext;
		if (!AC) return null;
		if (!ctx) ctx = new AC();
		if (ctx.state === 'suspended') ctx.resume();
		return ctx;
	}

	// Only real recordings are ever played: if a species has no file, nothing plays.
	function load(c, slug) {
		if (buffers[slug]) return Promise.resolve(buffers[slug]);
		if (!loading[slug]) {
			loading[slug] = fetch('/assets/audio/species/' + encodeURIComponent(slug) + '.mp3')
				.then((r) => (r.ok ? r.arrayBuffer() : null))
				.then((b) => (b ? c.decodeAudioData(b) : null))
				.then((buf) => { if (buf) buffers[slug] = buf; return buf; })
				.catch(() => null);
		}
		return loading[slug];
	}

	let current = null;
	function stopCurrent() {
		if (current) { try { current.stop(); } catch (e) {} current = null; }
	}

	// play returns a promise of how long the clip lasts in ms (0 if none).
	function play(slug) {
		const c = audio();
		if (!c || !slug) return Promise.resolve(0);
		return load(c, slug).then((buf) => {
			if (!buf) return 0;
			stopCurrent();
			const src = c.createBufferSource();
			src.buffer = buf;
			const gain = c.createGain();
			gain.gain.value = 1.4;
			src.connect(gain).connect(c.destination);
			src.loop = buf.duration < 4;
			src.start();
			let ms = buf.duration * 1000;
			if (src.loop) {
				ms = Math.max(ms * 3, 4000);
				src.stop(c.currentTime + ms / 1000);
			}
			current = src;
			return ms;
		});
	}

	document.addEventListener('click', (e) => {
		const btn = e.target.closest('[data-sound]');
		if (!btn) return;
		btn.classList.add('is-loading');
		play(btn.dataset.sound).then((ms) => {
			btn.classList.remove('is-loading');
			if (!ms) return;
			btn.classList.add('is-playing');
			clearTimeout(btn._t);
			btn._t = setTimeout(() => btn.classList.remove('is-playing'), ms);
		});
	});

	// Quizzes: hear a call, read a clue, or see a photo, then name the animal.
	const quiz = document.querySelector('[data-quiz]');
	if (quiz) {
		let animals = [];
		try { animals = JSON.parse(document.getElementById('quiz-data').textContent) || []; } catch (e) {}
		const stage = quiz.querySelector('[data-quiz-stage]');
		const status = quiz.querySelector('[data-quiz-status]');
		const live = quiz.querySelector('[data-quiz-live]');
		const tabs = [...quiz.querySelectorAll('[data-mode]')];
		const ROUNDS = 10;
		const shuffle = (a) => a.map((v) => [Math.random(), v]).sort((x, y) => x[0] - y[0]).map((x) => x[1]);
		const el = (tag, cls, text) => {
			const e = document.createElement(tag);
			if (cls) e.className = cls;
			if (text) e.textContent = text;
			return e;
		};
		const escapeRe = (s) => s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
		// Hide the animal's own name in its clue, or the answer is in the question.
		const mask = (fact, name) => name.split(/[\s-]+/).filter((w) => w.length > 2)
			.reduce((t, w) => t.replace(new RegExp('\\b' + escapeRe(w) + '(e?s)?\\b', 'gi'), '▒▒▒'), fact);

		const modes = {
			sound: { pool: () => animals.filter((a) => a.sound), cards: true },
			fact: { pool: () => animals.filter((a) => a.fact), cards: true },
			photo: { pool: () => animals, cards: false },
		};
		let mode = 'sound', order = [], n = 0, score = 0;

		function playButton(sound) {
			const b = el('button', 'hear quiz__play');
			b.type = 'button';
			const icon = el('span', 'hear__icon', '▶');
			icon.setAttribute('aria-hidden', 'true');
			b.append(icon, el('span', 'hear__label', 'Play the sound'));
			b.addEventListener('click', () => {
				b.classList.add('is-loading');
				play(sound).then((ms) => {
					b.classList.remove('is-loading');
					if (!ms) return;
					b.classList.add('is-playing');
					clearTimeout(b._t);
					b._t = setTimeout(() => b.classList.remove('is-playing'), ms);
				});
			});
			return b;
		}

		// Three wrong answers: look-alikes from the same group first, never the same genus.
		function distractors(current) {
			const others = animals.filter((a) => a !== current && (mode !== 'sound' || a.sound) && (mode !== 'fact' || a.fact) && a.genus !== current.genus);
			const same = shuffle(others.filter((a) => a.group === current.group));
			const rest = shuffle(others.filter((a) => a.group !== current.group));
			return (mode === 'photo' ? same.concat(rest) : shuffle(same.slice(0, 1).concat(rest))).slice(0, 3);
		}

		function status_() {
			status.textContent = 'Question ' + (n + 1) + ' of ' + order.length + ' · Score ' + score;
		}

		function finish() {
			stopCurrent();
			stage.textContent = '';
			status.textContent = '';
			const total = order.length;
			const msg = score === total ? 'Perfect! You really know your animals.' : score >= total * 0.7 ? 'Great work.' : score >= total / 2 ? 'Nice going.' : 'Keep exploring, it gets easier.';
			const box = el('div', 'quiz__result');
			box.append(el('p', 'quiz__score', score + ' / ' + total), el('p', 'quiz__msg', msg));
			const actions = el('p', 'quiz__actions');
			const again = el('button', 'btn-orange', 'Play again');
			again.type = 'button';
			again.addEventListener('click', begin);
			actions.append(again);
			tabs.filter((t) => t.dataset.mode !== mode).forEach((t) => {
				const b = el('button', 'button', 'Try “' + t.querySelector('.quizmode__title').textContent + '”');
				b.type = 'button';
				b.addEventListener('click', () => select(t.dataset.mode, true));
				actions.append(b);
			});
			box.append(actions);
			stage.append(box);
			live.textContent = 'Quiz finished. You got ' + score + ' of ' + total + '.';
			again.focus();
		}

		function ask() {
			stopCurrent();
			stage.textContent = '';
			if (n >= order.length) return finish();
			const current = order[n];
			status_();
			const options = shuffle([current, ...distractors(current)]);

			const prompt = el('div', 'quiz__prompt');
			let player = null;
			if (mode === 'sound') {
				player = playButton(current.sound);
				prompt.append(player);
			} else if (mode === 'fact') {
				const clue = el('blockquote', 'quiz__clue', mask(current.fact, current.name));
				prompt.append(clue);
			} else {
				const img = el('img', 'quiz__photo');
				img.src = current.photo;
				img.alt = 'Which animal is this?';
				prompt.append(img);
			}

			const grid = el('div', modes[mode].cards ? 'quiz__grid' : 'quiz__names');
			grid.setAttribute('role', 'group');
			grid.setAttribute('aria-label', 'Pick the animal');
			const feedback = el('div', 'quiz__feedback');

			options.forEach((o) => {
				const b = el('button', modes[mode].cards ? 'quiz__opt' : 'quiz__name');
				b.type = 'button';
				if (modes[mode].cards) {
					const img = el('img');
					img.src = o.photo;
					img.alt = '';
					img.loading = 'lazy';
					b.append(img);
				}
				b.append(el('span', '', o.name));
				b.addEventListener('click', () => {
					const right = o === current;
					if (right) score++;
					grid.querySelectorAll('button').forEach((x, i) => {
						x.disabled = true;
						if (options[i] === current) x.classList.add('is-right');
					});
					if (!right) b.classList.add('is-wrong');
					status_();

					const verdict = (right ? 'Correct! ' : 'Not quite. ') + 'It was the ' + current.name + '.';
					feedback.append(el('p', 'quiz__verdict ' + (right ? 'is-right' : 'is-wrong'), verdict));
					if (current.fact && mode !== 'fact') {
						const f = el('p', 'fact');
						const icon = el('span', 'fact__icon', '💡');
						icon.setAttribute('aria-hidden', 'true');
						const t = el('span');
						t.append(el('strong', 'fact__label', 'Did you know: '), document.createTextNode(current.fact));
						f.append(icon, t);
						feedback.append(f);
					}
					const actions = el('p', 'quiz__actions');
					const nextBtn = el('button', 'btn-orange', n + 1 < order.length ? 'Next question' : 'See your score');
					nextBtn.type = 'button';
					nextBtn.addEventListener('click', () => { n++; ask(); });
					const more = el('a', 'button', 'Meet the ' + current.name + ' →');
					more.href = current.href;
					actions.append(nextBtn, more);
					feedback.append(actions);
					live.textContent = verdict;
					nextBtn.focus();
				});
				grid.append(b);
			});
			stage.append(prompt, grid, feedback);
			if (status.getBoundingClientRect().top < 0 || prompt.getBoundingClientRect().bottom > window.innerHeight) {
				status.scrollIntoView({ block: 'start', behavior: 'smooth' });
			}
			if (player) { player.focus({ preventScroll: true }); player.click(); } else grid.querySelector('button').focus({ preventScroll: true });
		}

		function begin() {
			order = shuffle(modes[mode].pool()).slice(0, ROUNDS);
			n = 0;
			score = 0;
			ask();
		}

		function intro() {
			stopCurrent();
			stage.textContent = '';
			status.textContent = '';
			const pool = modes[mode].pool();
			if (pool.length < 4) {
				stage.append(el('p', '', 'This quiz needs at least four animals.'));
				return;
			}
			const start = el('button', 'btn-orange', 'Start · ' + Math.min(ROUNDS, pool.length) + ' questions');
			start.type = 'button';
			start.addEventListener('click', begin);
			stage.append(start);
		}

		function select(m, startNow) {
			mode = m;
			tabs.forEach((t) => t.setAttribute('aria-selected', String(t.dataset.mode === m)));
			if (startNow) begin(); else intro();
		}
		tabs.forEach((t) => t.addEventListener('click', () => select(t.dataset.mode)));
		select(new URLSearchParams(location.search).get('mode') in modes ? new URLSearchParams(location.search).get('mode') : 'sound');
	}
})();
