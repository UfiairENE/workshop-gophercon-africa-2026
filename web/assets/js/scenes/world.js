import { mount } from '/assets/js/threejs.js';
import { planetMaterial, dotTexture } from '/assets/js/scenes/planet-material.js';

// One body, large, for the top of its own page: a turning planet, or a
// star glowing in its own light.
mount('world', function ({ THREE, scene, camera, canvas, pointer }) {
	const look = canvas.parentElement.dataset;
	const star = 'star' in look;
	const ringed = 'ring' in look;
	camera.position.set(0, 0, 6);

	const group = new THREE.Group();
	group.rotation.z = 0.3;
	scene.add(group);

	const size = ringed ? 1.15 : 1.75;
	const body = new THREE.Mesh(
		new THREE.SphereGeometry(size, 96, 96),
		planetMaterial(THREE, {
			colour: look.colour,
			accent: look.accent,
			bands: parseFloat(look.bands) || 0,
			glow: star ? 0.45 : 0,
		}),
	);
	group.add(body);

	if (star) {
		const halo = new THREE.Sprite(
			new THREE.SpriteMaterial({ map: dotTexture(THREE, hexToRgba(look.colour, 0.75), hexToRgba(look.colour, 0)), depthWrite: false }),
		);
		halo.scale.set(size * 2.9, size * 2.9, 1);
		halo.renderOrder = -1;
		scene.add(halo);
	}

	if (ringed) {
		const rings = new THREE.Mesh(
			new THREE.RingGeometry(size * 1.3, size * 1.95, 128),
			new THREE.MeshBasicMaterial({ color: look.accent, side: THREE.DoubleSide, transparent: true, opacity: 0.55 }),
		);
		rings.rotation.x = 1.3;
		group.add(rings);
	}

	return function (time) {
		body.rotation.y = time * (star ? 0.05 : 0.12);
		body.material.uniforms.uTime.value = time;
		group.rotation.x += (pointer.y * 0.15 - group.rotation.x) * 0.05;
		group.rotation.y += (pointer.x * 0.25 - group.rotation.y) * 0.05;
	};
});

function hexToRgba(hex, alpha) {
	const n = parseInt(hex.slice(1), 16);
	return `rgba(${(n >> 16) & 255},${(n >> 8) & 255},${n & 255},${alpha})`;
}
