import { useRef, useEffect } from 'react';
import * as THREE from 'three';
import { ImprovedNoise } from 'three/examples/jsm/math/ImprovedNoise.js';

const LiquidMetalSphere = () => {
  const mountRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    if (!mountRef.current) return;
    const mouse = new THREE.Vector2();
    const scene = new THREE.Scene();
    const camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 1, 100);
    const renderer = new THREE.WebGLRenderer({ antialias: true });
    renderer.setSize(window.innerWidth, window.innerHeight);
    mountRef.current.appendChild(renderer.domElement);

    // Increase the radius to make the sphere larger
    const geometry = new THREE.SphereGeometry(4, 64, 64); // Radius increased to 3
    const vertexShader = `
    varying vec3 vNormal;
    varying vec3 vPosition;
    varying vec2 vUv;
    void main() {
      vNormal = normalize(normalMatrix * normal);
      vPosition = position;
      vec4 modelViewPosition = modelViewMatrix * vec4(position, 1.0);
      vUv = vec2(modelViewPosition.x, modelViewPosition.y) / 10.0 + 0.5;
      gl_Position = projectionMatrix * modelViewPosition;
    }
  `;
    const fragmentShader = `
    uniform float time;
    uniform vec2 mousePos;
    uniform vec2 resolution;
    varying vec3 vNormal;
    varying vec3 vPosition;
    varying vec2 vUv;
    void main() {
      vec3 orange = vec3(1.0, 0.55, 0.0);
      vec3 darkArea = vec3(0.1, 0.1, 0.1);
      float dist = distance(gl_PointCoord.xy, mousePos);
      float ripple = sin(dist * 10.0 - time * 5.0) * 0.5 + 0.5;
      vec3 color = mix(orange, darkArea, ripple);
      gl_FragColor = vec4(color, 64.0);
    }
  `;

    const material = new THREE.ShaderMaterial({
      uniforms: {
        time: { value: 1.0 },
        mousePos: { value: new THREE.Vector2(-1, -1) }, // Initialize outside the screen
        resolution: { value: new THREE.Vector2(window.innerWidth, window.innerHeight) },
      },
      vertexShader,
      fragmentShader,
    });

    const sphere = new THREE.Mesh(geometry, material);
    scene.add(sphere);
    camera.position.z = 10; // Adjust camera distance to fit the larger sphere

    const noise = new ImprovedNoise();
    const positionAttribute = geometry.getAttribute('position');
    const originalPosition: any[] = [];
    for (let i = 0; i < positionAttribute.count; i++) {
      originalPosition.push(new THREE.Vector3().fromBufferAttribute(positionAttribute, i));
    }
    const onMouseMove = (event: { clientX: number; clientY: number }) => {
      // Normalize mouse coordinates and update uniform
      material.uniforms.mousePos.value.set((event.clientX / window.innerWidth) * 2 - 1, -(event.clientY / window.innerHeight) * 2 + 1);
    };

    window.addEventListener('mousemove', onMouseMove);

    const animate = () => {
      requestAnimationFrame(animate);
      const time = Date.now() * 0.001;
      sphere.material.uniforms.time.value = time;

      for (let i = 0; i < positionAttribute.count; i++) {
        const vertex = originalPosition[i];
        const offset = noise.noise(vertex.x + time, vertex.y, vertex.z);
        const newPosition = vertex.clone().multiplyScalar(1 + offset * 0.03);
        positionAttribute.setXYZ(i, newPosition.x, newPosition.y, newPosition.z);
      }
      positionAttribute.needsUpdate = true;

      sphere.rotation.x += 0.005;
      sphere.rotation.y += 0.005;
      renderer.render(scene, camera);
    };

    animate();

    return () => {
      renderer.dispose();
      window.removeEventListener('mousemove', onMouseMove);
      window.removeEventListener('resize', () => {});
    };
  }, []);

  return <div ref={mountRef} />;
};

export default LiquidMetalSphere;
