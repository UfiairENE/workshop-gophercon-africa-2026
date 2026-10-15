// A flat-shaded, banded planet surface in the site's pastel palette,
// shared by the orrery and the portraits.
export function planetMaterial(THREE, { colour, accent, bands = 0.3, glow = 0 }) {
	return new THREE.ShaderMaterial({
		uniforms: {
			uColour: { value: new THREE.Color(colour) },
			uAccent: { value: new THREE.Color(accent || colour) },
			uBands: { value: bands },
			uGlow: { value: glow },
			uTime: { value: 0 },
			uLight: { value: new THREE.Vector3(-0.6, 0.55, 0.75).normalize() },
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
			uniform vec3 uColour;
			uniform vec3 uAccent;
			uniform float uBands;
			uniform float uGlow;
			uniform float uTime;
			uniform vec3 uLight;
			varying vec3 vNormal;
			varying vec3 vView;
			varying vec2 vUv;
			void main() {
				float wobble = sin(vUv.x * 12.566 + uTime * 0.2) * 0.012;
				float stripes = sin((vUv.y + wobble) * 28.0) * 0.5 + 0.5;
				stripes *= sin((vUv.y + wobble) * 9.0 + 1.3) * 0.5 + 0.5;
				vec3 base = mix(uColour, uAccent, stripes * uBands);
				float lit = smoothstep(-0.3, 1.0, dot(vNormal, uLight));
				vec3 shade = mix(base * vec3(0.62, 0.62, 0.78), base, lit);
				float rim = pow(1.0 - max(dot(vNormal, vView), 0.0), 2.5);
				vec3 colour = shade + vec3(1.0) * rim * 0.18;
				colour = mix(colour, uColour * 1.08, uGlow);
				gl_FragColor = vec4(colour, 1.0);
			}
		`,
	});
}

// A soft round dot, for glows and star sprites.
export function dotTexture(THREE, inner = 'rgba(255,255,255,1)', outer = 'rgba(255,255,255,0)') {
	const size = 128;
	const canvas = document.createElement('canvas');
	canvas.width = canvas.height = size;
	const ctx = canvas.getContext('2d');
	const g = ctx.createRadialGradient(size / 2, size / 2, 0, size / 2, size / 2, size / 2);
	g.addColorStop(0, inner);
	g.addColorStop(0.35, inner);
	g.addColorStop(1, outer);
	ctx.fillStyle = g;
	ctx.fillRect(0, 0, size, size);
	const texture = new THREE.CanvasTexture(canvas);
	texture.colorSpace = THREE.SRGBColorSpace;
	return texture;
}
