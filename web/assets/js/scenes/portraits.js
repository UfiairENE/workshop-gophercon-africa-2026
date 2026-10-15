import { mount } from '/assets/js/threejs.js';
import { planetMaterial } from '/assets/js/scenes/planet-material.js';

// One canvas behind the whole grid draws every portrait into its own
// box, so a dozen planets cost a single WebGL context.
mount('portraits', function ({ THREE, renderer, canvas, still }) {
	const views = Array.from(canvas.parentElement.querySelectorAll('.portrait__view'));
	const camera = new THREE.PerspectiveCamera(30, 1, 0.1, 50);
	camera.position.set(0, 0, 5.4);

	const portraits = views.map(function (el, i) {
		const scene = new THREE.Scene();
		const material = planetMaterial(THREE, {
			colour: el.dataset.colour,
			accent: el.dataset.accent,
			bands: parseFloat(el.dataset.bands) || 0,
		});
		const ringed = el.hasAttribute('data-ring');
		const planet = new THREE.Mesh(new THREE.SphereGeometry(ringed ? 0.72 : 1, 64, 64), material);
		planet.rotation.z = 0.25 + (i % 3) * 0.12;
		scene.add(planet);

		if (ringed) {
			const rings = new THREE.Mesh(
				new THREE.RingGeometry(0.92, 1.36, 96),
				new THREE.MeshBasicMaterial({ color: el.dataset.accent, side: THREE.DoubleSide, transparent: true, opacity: 0.6 }),
			);
			rings.rotation.x = 1.35;
			rings.rotation.y = -0.35;
			scene.add(rings);
		}

		return { el, card: el.closest('.portrait'), scene, planet, material, speed: 0.25, scale: 1, angle: i * 1.3 };
	});

	return function (time, delta) {
		const box = canvas.getBoundingClientRect();
		renderer.setScissorTest(false);
		renderer.clear();
		renderer.setScissorTest(true);

		for (const p of portraits) {
			const r = p.el.getBoundingClientRect();
			if (r.width === 0 || (!still && (r.bottom < 0 || r.top > window.innerHeight))) continue;

			// Hovering a card spins its planet up and draws it closer.
			const hovered = !still && p.card.matches(':hover, :focus-within');
			p.speed += ((hovered ? 2.2 : 0.25) - p.speed) * 0.08;
			p.scale += ((hovered ? 1.1 : 1) - p.scale) * 0.12;
			p.angle += p.speed * (delta || 0);

			p.planet.rotation.y = p.angle;
			p.planet.scale.setScalar(p.scale);
			p.material.uniforms.uTime.value = time;

			const x = r.left - box.left;
			const y = box.bottom - r.bottom;
			renderer.setViewport(x, y, r.width, r.height);
			renderer.setScissor(x, y, r.width, r.height);
			camera.aspect = r.width / r.height;
			camera.updateProjectionMatrix();
			renderer.render(p.scene, camera);
		}
		return false;
	};
});
