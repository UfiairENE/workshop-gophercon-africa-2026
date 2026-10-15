import { mount } from '/assets/js/threejs.js';

mount('planet', function ({ THREE, scene, camera, pointer, colour }) {
	camera.position.set(0, 0, 6.2);

	const planet = new THREE.Mesh(
		new THREE.SphereGeometry(1.75, 96, 96),
		new THREE.ShaderMaterial({
			uniforms: {
				uWarm: { value: colour('--three-warm', '#ffcf9e') },
				uCool: { value: colour('--three-cool', '#babfff') },
				uBloom: { value: colour('--three-bloom', '#e3a2ef') },
				uSky: { value: colour('--three-sky', '#bfe0f7') },
				uTime: { value: 0 },
			},
			vertexShader: `
				varying vec3 vNormal;
				varying vec3 vView;
				varying vec2 vUv;
				void main() {
					vUv = uv;
					vNormal = normalize(normalMatrix * normal);
					vec4 mv = modelViewMatrix * vec4(position, 1.0);
					vView = normalize(-mv.xyz);
					gl_Position = projectionMatrix * mv;
				}
			`,
			fragmentShader: `
				uniform vec3 uWarm;
				uniform vec3 uCool;
				uniform vec3 uBloom;
				uniform vec3 uSky;
				uniform float uTime;
				varying vec3 vNormal;
				varying vec3 vView;
				varying vec2 vUv;
				void main() {
					float lat = vUv.y;
					float drift = sin(vUv.x * 6.2832 * 2.0 + uTime * 0.15) * 0.04;
					float bands = sin((lat + drift) * 22.0) * 0.5 + 0.5;
					vec3 base = mix(uCool, uWarm, smoothstep(0.25, 0.75, lat));
					base = mix(base, uBloom, bands * 0.28);
					base = mix(base, uSky, smoothstep(0.85, 1.0, lat) * 0.6);
					vec3 light = normalize(vec3(-0.6, 0.7, 0.8));
					float lit = smoothstep(-0.35, 1.0, dot(vNormal, light));
					vec3 shade = mix(uCool * 0.55, base, lit);
					float rim = pow(1.0 - max(dot(vNormal, vView), 0.0), 3.0);
					gl_FragColor = vec4(shade + uWarm * rim * 0.35, 1.0);
				}
			`,
		}),
	);
	planet.rotation.z = 0.35;
	scene.add(planet);

	const orbitTilt = new THREE.Group();
	orbitTilt.rotation.set(1.22, 0, 0.32);
	scene.add(orbitTilt);

	const points = [];
	for (let i = 0; i <= 256; i++) {
		const a = (i / 256) * Math.PI * 2;
		points.push(new THREE.Vector3(Math.cos(a) * 2.7, Math.sin(a) * 2.7, 0));
	}
	orbitTilt.add(
		new THREE.Line(
			new THREE.BufferGeometry().setFromPoints(points),
			new THREE.LineBasicMaterial({ color: colour('--three-line', '#3a3a3e'), transparent: true, opacity: 0.45 }),
		),
	);

	const moon = new THREE.Mesh(
		new THREE.SphereGeometry(0.16, 32, 32),
		new THREE.MeshBasicMaterial({ color: colour('--three-line', '#3a3a3e') }),
	);
	orbitTilt.add(moon);

	return function (time) {
		planet.material.uniforms.uTime.value = time;
		planet.rotation.y = time * 0.08;
		const a = time * 0.25 + 0.8;
		moon.position.set(Math.cos(a) * 2.7, Math.sin(a) * 2.7, 0);
		scene.rotation.x += (pointer.y * 0.08 - scene.rotation.x) * 0.04;
		scene.rotation.y += (pointer.x * 0.12 - scene.rotation.y) * 0.04;
	};
});
