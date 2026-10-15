// Draws every [data-map] on the page with Leaflet, reading its points from
// the JSON the component wrote out. Markers always go on a canvas, since
// the home page has thousands; clusters and heat are opt-in per map and
// used only when their plugin loaded. Without Leaflet the fallback shows.

(function () {
	const TILES = 'https://server.arcgisonline.com/ArcGIS/rest/services/NatGeo_World_Map/MapServer/tile/{z}/{y}/{x}';
	const CREDIT = 'Tiles &copy; Esri &mdash; National Geographic, Esri, DeLorme, NAVTEQ, UNEP-WCMC, USGS, NASA, ESA, METI, NRCAN, GEBCO, NOAA, iPC';
	const still = window.matchMedia('(prefers-reduced-motion: reduce)').matches;

	function points(id) {
		const el = document.getElementById(id);
		try {
			return (el && JSON.parse(el.textContent)) || [];
		} catch (err) {
			return [];
		}
	}

	// Built from nodes, never an HTML string: every field here is GBIF's text.
	function popup(point) {
		const root = document.createElement('div');
		root.className = 'map__popup';
		if (point.photo) {
			const img = root.appendChild(document.createElement('img'));
			img.className = 'map__photo';
			img.src = point.photo;
			img.alt = '';
			img.loading = 'lazy';
		}
		root.appendChild(document.createElement('strong')).textContent = point.title;
		root.appendChild(document.createElement('span')).textContent =
			[point.date, point.place, point.count].filter(Boolean).join(' · ');
		if (point.href) {
			const link = root.appendChild(document.createElement('a'));
			link.href = point.href;
			link.textContent = 'Open sighting';
		}
		return root;
	}

	// Colours come from the element's CSS, so the palette lives in one place.
	function colour(el, prop, fallback) {
		return getComputedStyle(el).getPropertyValue(prop).trim() || fallback;
	}

	function init(el) {
		const data = points(el.dataset.map);
		if (!data.length || typeof L === 'undefined') {
			el.classList.add('is-unavailable');
			return;
		}

		const map = L.map(el, {
			preferCanvas: true,
			scrollWheelZoom: false,
			zoomAnimation: !still,
			fadeAnimation: !still,
			markerZoomAnimation: !still,
			worldCopyJump: true,
		});
		L.tileLayer(el.dataset.tiles || TILES, {
			attribution: el.dataset.attribution || CREDIT,
			maxZoom: 16,
			minZoom: 2,
		}).addTo(map);

		const latLngs = data.map((p) => [p.lat, p.lng]);
		let layer;
		if ('heat' in el.dataset && typeof L.heatLayer === 'function') {
			layer = L.heatLayer(latLngs, { radius: 18, blur: 15, minOpacity: 0.3 });
		} else {
			const fill = colour(el, '--map-marker', '#c2410c');
			const stroke = colour(el, '--map-marker-stroke', '#fff');
			layer = 'cluster' in el.dataset && typeof L.markerClusterGroup === 'function'
				? L.markerClusterGroup({ chunkedLoading: true, showCoverageOnHover: false, maxClusterRadius: 40 })
				: L.featureGroup();
			const markers = data.map((p) => ({
				group: p.group,
				marker: L.circleMarker([p.lat, p.lng], { radius: 7, weight: 2, color: stroke, fillColor: p.colour || fill, fillOpacity: 0.9 })
					// A function, so the popup is only built when it opens.
					.bindPopup(() => popup(p), { minWidth: 220 })
					// The animal's name on hover; always shown when the map has one sighting.
					.bindTooltip(p.title, { permanent: data.length === 1, direction: 'top', offset: [0, -8], className: 'map__label' }),
			}));
			markers.forEach((m) => m.marker.addTo(layer));

			// The legend's buttons show and hide each animal group.
			const legend = document.querySelector('[data-map-legend="' + el.dataset.map + '"]');
			if (legend) {
				legend.addEventListener('click', (e) => {
					const btn = e.target.closest('[data-group]');
					if (!btn) return;
					btn.setAttribute('aria-pressed', btn.getAttribute('aria-pressed') === 'true' ? 'false' : 'true');
					const on = new Set([...legend.querySelectorAll('[aria-pressed="true"]')].map((b) => b.dataset.group));
					layer.clearLayers();
					markers.filter((m) => on.has(m.group)).forEach((m) => m.marker.addTo(layer));
				});
			}
		}
		layer.addTo(map);

		map.fitBounds(L.latLngBounds(latLngs), {
			padding: [24, 24],
			animate: false,
			maxZoom: Number(el.dataset.maxZoom) || 12,
		});

		// Wheel zoom waits for a click or focus, so the page scrolls past freely.
		map.on('focus', () => map.scrollWheelZoom.enable());
		map.on('blur', () => map.scrollWheelZoom.disable());
	}

	document.querySelectorAll('[data-map]').forEach(init);
})();
