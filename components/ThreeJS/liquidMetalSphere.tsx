import { useRef, useEffect } from 'react';
import * as THREE from 'three';

const LiquidMetalSphere = () => {
  const mountRef = useRef<HTMLDivElement | null>(
    null,
  );

  useEffect(() => {
    if (!mountRef.current) return;

    const scene = new THREE.Scene();
    const camera = new THREE.PerspectiveCamera(
      75,
      window.innerWidth / window.innerHeight,
      1,
      10,
    );
    const renderer = new THREE.WebGLRenderer({
      antialias: true,
    });

    renderer.setSize(
      window.innerWidth,
      window.innerHeight,
    );
    mountRef.current.appendChild(
      renderer.domElement,
    );


    const geometry = new THREE.SphereGeometry(
      1,
      32,
      32,
    );
    const material = new THREE.MeshBasicMaterial({
      color: 'orange',
      wireframe: true,
    }); 
    const sphere = new THREE.Mesh(
      geometry,
      material,
    );
    scene.add(sphere);

    camera.position.z = 5;

    const animate = () => {
      requestAnimationFrame(animate);
      sphere.rotation.x += 0.001;
      sphere.rotation.y += 0.001;
      renderer.render(scene, camera);
    };

    animate();

    return () => {
      renderer.dispose();
      window.removeEventListener(
        'resize',
        () => {},
      );
    };
  }, []);

  return <div ref={mountRef} />;
};

export default LiquidMetalSphere;
