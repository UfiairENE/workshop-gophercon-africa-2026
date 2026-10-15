(function () {
	const menu = document.querySelector('[data-menu]');
	if (menu) {
		const nav = document.getElementById(menu.getAttribute('aria-controls'));
		const set = (open) => {
			menu.setAttribute('aria-expanded', String(open));
			document.body.classList.toggle('menu-open', open);
		};
		menu.addEventListener('click', () => set(menu.getAttribute('aria-expanded') !== 'true'));
		nav.addEventListener('click', (e) => { if (e.target.closest('a')) set(false); });
		document.addEventListener('keydown', (e) => {
			if (e.key === 'Escape' && menu.getAttribute('aria-expanded') === 'true') { set(false); menu.focus(); }
		});
	}

	const box = document.getElementById('lightbox');
	document.querySelectorAll('[data-lightbox]').forEach((t) => {
		t.addEventListener('click', () => {
			if (!box || !box.showModal) {
				window.open(t.dataset.lightbox, '_blank', 'noopener');
				return;
			}
			box.querySelector('img').src = t.dataset.lightbox;
			box.querySelector('img').alt = t.dataset.caption || '';
			box.querySelector('figcaption span').textContent = t.dataset.caption || '';
			const a = box.querySelector('figcaption a');
			if (t.dataset.href) { a.href = t.dataset.href; a.hidden = false; } else { a.hidden = true; }
			box.showModal();
		});
	});
	if (box) {
		box.addEventListener('click', (e) => {
			if (e.target === box || e.target.closest('.lightbox__close')) box.close();
		});
	}

	// A real hero video, if one has been dropped in.
	const heroBox = document.querySelector('.banner__image');
	if (heroBox && !window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
		fetch('/assets/video/hero.mp4', { method: 'HEAD' }).then((r) => {
			if (!r.ok || !/video/.test(r.headers.get('content-type') || '')) return;
			// Phones only autoplay a video that is muted and inline before it has a
			// source, and don't buffer until asked to play, so play straight away
			// and show it once frames are actually moving.
			const v = document.createElement('video');
			v.className = 'banner__video';
			v.muted = true;
			v.defaultMuted = true;
			v.playsInline = true;
			v.loop = true;
			v.autoplay = true;
			['muted', 'playsinline', 'autoplay', 'loop'].forEach((a) => v.setAttribute(a, ''));
			v.setAttribute('aria-hidden', 'true');
			v.preload = 'auto';
			v.src = '/assets/video/hero.mp4';
			v.addEventListener('playing', () => heroBox.classList.add('has-video'), { once: true });
			heroBox.appendChild(v);
			// Low Power Mode and data savers block autoplay; a touch lets it start.
			v.play().catch(() => {
				const retry = () => v.play().catch(() => {});
				['touchend', 'click', 'keydown'].forEach((t) => window.addEventListener(t, retry, { once: true, passive: true }));
			});
		}).catch(() => {});
	}

	// Near me: nearest sightings from the points on the page.
	const near = document.querySelector('[data-near]');
	if (near) {
		let pts = [];
		try { pts = JSON.parse(document.getElementById('near-points').textContent) || []; } catch (e) {}
		const out = near.querySelector('[data-near-results]');
		const msg = near.querySelector('[data-near-status]');
		const km = (a, b, c, d) => {
			const R = 6371, r = Math.PI / 180;
			const x = Math.sin((c - a) * r / 2) ** 2 + Math.cos(a * r) * Math.cos(c * r) * Math.sin((d - b) * r / 2) ** 2;
			return 2 * R * Math.asin(Math.sqrt(x));
		};
		near.querySelector('[data-near-go]').addEventListener('click', () => {
			if (!navigator.geolocation) { msg.textContent = 'Your browser can’t share its location.'; return; }
			msg.textContent = 'Finding you…';
			navigator.geolocation.getCurrentPosition((pos) => {
				const { latitude: la, longitude: lo } = pos.coords;
				const ranked = pts.map((p) => ({ p, d: km(la, lo, p.lat, p.lng) })).sort((a, b) => a.d - b.d).slice(0, 12);
				out.textContent = '';
				ranked.forEach(({ p, d }) => {
					const li = document.createElement('li');
					const a = document.createElement('a');
					a.href = p.href;
					a.textContent = p.title + ', ' + p.place;
					const s = document.createElement('span');
					s.textContent = Math.round(d).toLocaleString('en-GB') + ' km away · ' + p.date;
					li.append(a, s);
					out.appendChild(li);
				});
				msg.textContent = ranked.length ? 'The ' + ranked.length + ' sightings nearest to you:' : 'No sightings with a position yet.';
			}, () => { msg.textContent = 'Location was blocked. Allow it in the address bar and try again.'; }, { timeout: 10000 });
		});
	}

	// Poster: a downloadable card for one sighting.
	document.querySelectorAll('[data-poster]').forEach((btn) => {
		btn.addEventListener('click', () => {
			const d = btn.dataset;
			const W = 1080, H = 1350;
			const cv = document.createElement('canvas');
			cv.width = W; cv.height = H;
			const g = cv.getContext('2d');

			function text(t, x, y, font, color, max) {
				g.font = font; g.fillStyle = color; g.textBaseline = 'alphabetic';
				let line = '', yy = y;
				t.split(' ').forEach((w) => {
					const test = line ? line + ' ' + w : w;
					if (max && g.measureText(test).width > max && line) { g.fillText(line, x, yy); line = w; yy += parseInt(font.match(/(\d+)px/)[1], 10) * 1.15; } else line = test;
				});
				g.fillText(line, x, yy);
				return yy;
			}

			function finish(img) {
				g.fillStyle = '#322c2a';
				g.fillRect(0, 0, W, H);
				if (img) {
					const r = Math.max(W / img.width, 860 / img.height);
					const w = img.width * r, h = img.height * r;
					g.drawImage(img, (W - w) / 2, (860 - h) / 2, w, h);
				}
				const grad = g.createLinearGradient(0, 520, 0, 860);
				grad.addColorStop(0, 'rgba(50,44,42,0)');
				grad.addColorStop(1, 'rgba(50,44,42,1)');
				g.fillStyle = grad;
				g.fillRect(0, 520, W, 342);
				let y = text(d.name, 72, 940, '300 84px Raleway, sans-serif', '#ffffff', W - 144);
				if (d.sci) y = text(d.sci, 72, y + 52, 'italic 400 38px "Crimson Pro", serif', '#e8d9cf', W - 144);
				y = text(d.place, 72, y + 64, '600 40px "Source Sans 3", sans-serif', '#ffffff', W - 144);
				y = text(d.date + '  ·  ' + d.temp + '  ·  ' + d.cond, 72, y + 52, '400 34px "Source Sans 3", sans-serif', '#d5d3ce', W - 144);
				g.fillStyle = 'rgb(233,114,0)';
				g.fillRect(0, H - 110, W, 110);
				text('FIELDLIGHT', 72, H - 42, '700 40px Raleway, sans-serif', '#ffffff');
				text('real sightings, real weather', W - 72 - 430, H - 42, '400 30px "Source Sans 3", sans-serif', '#ffffff');
				try {
					cv.toBlob((b) => {
						if (!b) return;
						const a = document.createElement('a');
						a.href = URL.createObjectURL(b);
						a.download = (d.name || 'sighting').toLowerCase().replace(/[^a-z0-9]+/g, '-') + '-poster.png';
						document.body.appendChild(a);
						a.click();
						a.remove();
					}, 'image/png');
				} catch (e) { btn.textContent = 'Poster unavailable'; }
			}

			if (!d.photo) { finish(null); return; }
			const img = new Image();
			img.crossOrigin = 'anonymous';
			img.onload = () => {
				try { finish(img); } catch (e) { finish(null); }
			};
			img.onerror = () => finish(null);
			img.src = d.photo;
		});
	});
})();
