import { gsap } from 'gsap';
import ScrollSmoother from 'gsap/dist/ScrollSmoother';
import { ScrollTrigger } from 'gsap/dist/ScrollTrigger';

if (typeof window !== "undefined") {
    gsap.registerPlugin(ScrollTrigger, ScrollSmoother);
}

export * from "gsap";
export { gsap, ScrollTrigger, ScrollSmoother };
