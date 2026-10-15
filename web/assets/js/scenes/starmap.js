import { mount, scrollProgress } from '/assets/js/threejs.js';
import { OrbitControls } from 'three/addons/controls/OrbitControls.js';
import { dotTexture } from '/assets/js/scenes/planet-material.js';

// Distance is squashed so the nearest stars and the far Kepler field both fit.
const scale = (parsecs) => 1.4 * Math.log10(1 + parsecs);

mount('starmap', function ({ THREE, scene, camera, canvas, still }) {
	const data = JSON.parse(document.getElementById('starmap-data').textContent || '{}');
	const systems = data.points || [];
	const names = data.names || [];
	const section = canvas.closest('[data-scroll-scene]');

	const hot = new THREE.Color('#bfe0f7');
	const sunlike = new THREE.Color('#ffcf9e');
	const cool = new THREE.Color('#e3a2ef');
	const unknown = new THREE.Color('#b5b0b0');

	const positions = new Float32Array(systems.length * 3);
	const colours = new Float32Array(systems.length * 3);
	systems.forEach(function ([ra, dec, pc, teff], i) {
		const a = THREE.MathUtils.degToRad(ra);
		const d = THREE.MathUtils.degToRad(dec);
		const r = scale(pc);
		positions[i * 3] = r * Math.cos(d) * Math.cos(a);
		positions[i * 3 + 1] = r * Math.sin(d);
		positions[i * 3 + 2] = -r * Math.cos(d) * Math.sin(a);
		const c = teff === 0 ? unknown : teff > 7000 ? hot : teff > 5000 ? sunlike : cool;
		colours.set([c.r, c.g, c.b], i * 3);
	});

	const geometry = new THREE.BufferGeometry();
	geometry.setAttribute('position', new THREE.BufferAttribute(positions, 3));
	geometry.setAttribute('color', new THREE.BufferAttribute(colours, 3));

	const field = new THREE.Group();
	scene.add(field);
	const stars = new THREE.Points(
			geometry,
			new THREE.PointsMaterial({
				size: 0.055,
				map: dotTexture(THREE),
				vertexColors: true,
				transparent: true,
				opacity: 0.8,
				depthWrite: false,
			}),
		);
	field.add(stars);

	const sun = new THREE.Sprite(
		new THREE.SpriteMaterial({ map: dotTexture(THREE, 'rgba(255,207,158,1)', 'rgba(255,207,158,0)'), depthWrite: false, blending: THREE.AdditiveBlending }),
	);
	sun.scale.set(0.45, 0.45, 1);
	field.add(sun);

	// Rings at 10, 100 and 1,000 parsecs, flat on the celestial equator.
	for (const pc of [10, 100, 1000]) {
		const r = scale(pc);
		const ring = [];
		for (let k = 0; k <= 160; k++) {
			const a = (k / 160) * Math.PI * 2;
			ring.push(new THREE.Vector3(Math.cos(a) * r, 0, Math.sin(a) * r));
		}
		field.add(
			new THREE.Line(
				new THREE.BufferGeometry().setFromPoints(ring),
				new THREE.LineBasicMaterial({ color: '#ffffff', transparent: true, opacity: 0.08 }),
			),
		);
	}

	camera.position.set(6, 3.5, 8);

	// Dragging on a phone would trap the page's scroll, so touch screens
	// get the scroll-driven flight only.
	const touch = window.matchMedia('(pointer: coarse)').matches;
	let controls = null;
	if (!touch) {
		controls = new OrbitControls(camera, canvas);
		controls.enableZoom = false;
		controls.enablePan = false;
		controls.enableDamping = true;
		controls.autoRotate = !still;
		controls.autoRotateSpeed = 0.35;
	} else {
		canvas.style.pointerEvents = 'none';
	}

	// Hovering names a star, and clicking opens its page.
	const tip = canvas.parentElement.querySelector('.starmap__tip');
	const raycaster = new THREE.Raycaster();
	raycaster.params.Points.threshold = 0.08;
	const mouse = new THREE.Vector2();
	let hovered = -1;
	let down = null;

	function pick(e) {
		const rect = canvas.getBoundingClientRect();
		mouse.set(((e.clientX - rect.left) / rect.width) * 2 - 1, -((e.clientY - rect.top) / rect.height) * 2 + 1);
		raycaster.setFromCamera(mouse, camera);
		const hit = raycaster.intersectObject(stars)[0];
		return hit ? hit.index : -1;
	}

	if (!touch) {
		canvas.addEventListener('pointermove', function (e) {
			hovered = pick(e);
			if (hovered < 0) {
				tip.hidden = true;
				canvas.style.cursor = '';
				return;
			}
			const rect = canvas.getBoundingClientRect();
			tip.textContent = names[hovered];
			tip.style.left = e.clientX - rect.left + 14 + 'px';
			tip.style.top = e.clientY - rect.top + 14 + 'px';
			tip.hidden = false;
			canvas.style.cursor = 'pointer';
		});
		canvas.addEventListener('pointerleave', function () {
			tip.hidden = true;
		});
		canvas.addEventListener('pointerdown', function (e) {
			down = [e.clientX, e.clientY];
		});
		canvas.addEventListener('pointerup', function (e) {
			// A drag turns the map; only a click that stayed put opens a star.
			if (!down || Math.hypot(e.clientX - down[0], e.clientY - down[1]) > 4) return;
			const i = pick(e);
			if (i >= 0) window.location.href = '/worlds/' + encodeURIComponent(names[i]);
		});
	}

	let distance = 11;

	return function (time) {
		// Scrolling through the section flies the camera in towards the Sun.
		const t = section ? scrollProgress(section) : 0.5;
		const target = 12 - 7.5 * THREE.MathUtils.smoothstep(t, 0.15, 0.85);
		distance += (target - distance) * 0.06;
		camera.position.setLength(distance * Math.max(1, 1.3 / camera.aspect));

		if (controls) {
			controls.update();
		} else {
			field.rotation.y = time * 0.04;
			camera.lookAt(0, 0, 0);
		}
	};
});
