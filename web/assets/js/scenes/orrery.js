import { mount } from '/assets/js/threejs.js';
import { planetMaterial, dotTexture } from '/assets/js/scenes/planet-material.js';

// One Earth year takes this many seconds on screen.
const SECONDS_PER_YEAR = 14;

mount('orrery', function ({ THREE, scene, camera, canvas, pointer, colour }) {
	const data = JSON.parse(document.getElementById('orrery-data').textContent || '[]');
	const section = canvas.closest('[data-scroll-scene]');
	const ink = colour('--three-line', '#3a3a3e');

	const orbitRadius = (au) => 0.8 + 1.6 * Math.pow(au, 0.4);
	const bodySize = (earths) => 0.07 + 0.09 * Math.log2(1 + earths);

	const system = new THREE.Group();
	scene.add(system);

	const sun = new THREE.Mesh(
		new THREE.SphereGeometry(0.55, 64, 64),
		planetMaterial(THREE, { colour: '#ffcf9e', accent: '#ffb877', bands: 0.15, glow: 0.55 }),
	);
	system.add(sun);

	const halo = new THREE.Sprite(
		new THREE.SpriteMaterial({ map: dotTexture(THREE, 'rgba(255,207,158,0.55)', 'rgba(255,207,158,0)'), depthWrite: false }),
	);
	halo.scale.set(3.4, 3.4, 1);
	system.add(halo);

	const planets = data.map(function (b, i) {
		const r = orbitRadius(b.a);
		const ring = [];
		for (let k = 0; k <= 192; k++) {
			const a = (k / 192) * Math.PI * 2;
			ring.push(new THREE.Vector3(Math.cos(a) * r, 0, Math.sin(a) * r));
		}
		system.add(
			new THREE.Line(
				new THREE.BufferGeometry().setFromPoints(ring),
				new THREE.LineBasicMaterial({ color: ink, transparent: true, opacity: b.kind === 'planet' ? 0.28 : 0.12 }),
			),
		);

		const size = bodySize(b.radius);
		const mesh = new THREE.Mesh(new THREE.SphereGeometry(size, 48, 48), planetMaterial(THREE, { colour: b.colour, bands: b.kind === 'planet' ? 0.4 : 0.1 }));
		if (b.ring) {
			const rings = new THREE.Mesh(
				new THREE.RingGeometry(size * 1.35, size * (b.name === 'Saturn' ? 1.9 : 1.6), 64),
				new THREE.MeshBasicMaterial({ color: b.name === 'Saturn' ? '#ffcf9e' : '#bfe0f7', side: THREE.DoubleSide, transparent: true, opacity: 0.45 }),
			);
			rings.rotation.x = b.name === 'Saturn' ? 1.35 : 0.2;
			mesh.add(rings);
		}
		system.add(mesh);

		return { mesh, r, rate: (Math.PI * 2 * 365.256) / (b.period * SECONDS_PER_YEAR), phase: i * 2.17 };
	});

	// Clicking a planet opens its page.
	const raycaster = new THREE.Raycaster();
	const mouse = new THREE.Vector2();
	const meshes = planets.map((p) => p.mesh);

	function pick(e) {
		const rect = canvas.getBoundingClientRect();
		mouse.set(((e.clientX - rect.left) / rect.width) * 2 - 1, -((e.clientY - rect.top) / rect.height) * 2 + 1);
		raycaster.setFromCamera(mouse, camera);
		const hit = raycaster.intersectObjects([sun, ...meshes], false)[0];
		if (!hit) return null;
		return hit.object === sun ? 'Sun' : data[meshes.indexOf(hit.object)].name;
	}

	canvas.addEventListener('pointermove', function (e) {
		canvas.style.cursor = pick(e) ? 'pointer' : '';
	});
	canvas.addEventListener('click', function (e) {
		const name = pick(e);
		if (name) window.location.href = '/worlds/' + encodeURIComponent(name);
	});

	let tilt = 1.05;
	let distance = 15;

	return function (time) {
		// Scrolling down through the banner swings the view from overhead
		// to edge-on and draws the camera in.
		let t = 0;
		if (section) {
			const rect = section.getBoundingClientRect();
			t = Math.min(1, Math.max(0, -rect.top / rect.height));
		}
		const narrow = Math.max(1, 1.5 / camera.aspect);
		tilt += (1.05 - 0.85 * t + pointer.y * 0.06 - tilt) * 0.08;
		distance += ((15 - 5 * t) * narrow - distance) * 0.08;

		const spin = time * 0.02 + pointer.x * 0.15;
		camera.position.set(Math.sin(spin) * Math.cos(tilt) * distance, Math.sin(tilt) * distance, Math.cos(spin) * Math.cos(tilt) * distance);
		camera.lookAt(0, -0.4, 0);

		sun.material.uniforms.uTime.value = time;
		sun.rotation.y = time * 0.05;
		for (const p of planets) {
			const a = p.phase + time * p.rate;
			p.mesh.position.set(Math.cos(a) * p.r, 0, -Math.sin(a) * p.r);
			p.mesh.rotation.y = time * 0.6;
			p.mesh.material.uniforms.uTime.value = time;
		}
	};
});
