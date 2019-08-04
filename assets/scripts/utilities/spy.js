import debounce from './debouncer';

/**
 * Spy sections and add class when in view
 * @param {string} sectionHandle
 * @param {string} inViewClass
 */
export function spySections(sectionHandle = '.spy', inViewClass = 'in-view') {
	const sections = document.querySelectorAll(sectionHandle);

	// Event handler functions
	function handleScroll(entries, observer) {
		if (entries.length === 0) {
			observer.disconnect();
		}

		Array.from(entries)
			.filter(entry => (
				entry.target.classList.contains(inViewClass)
                || entry.intersectionRatio > (entry.target.getAttribute('data-threshold') || 0.5)
			))
			.forEach(entry => {
				entry.target.classList.add(inViewClass);
				observer.unobserve(entry.target);
			});
	}

	// Initialize
	const observer = new IntersectionObserver(handleScroll, {
		// Call at multiple thresholds to allow for customization via data attribute
		threshold: [0, 0.25, 0.5, 0.75, 1],
	});
	Array.from(sections).forEach(section => { observer.observe(section); });
}

/**
 * Spy section and set correspond link to active
 * @param {object}
 * |_@param {nodeList} links
 * |_@param {nodeList} links
 * |_@param {string} activeClass
 * |_@param {float} threshold
 * |_@param {function} cb
 */
export function spyAnchors({
	links = [],
	sections = [],
	activeClass = 'is-active',
	threshold = window.innerHeight / 2,
	scrollContainer = window,
	cb = null,
}) {
	if (!sections) return;

	function handleScroll() {
		Array.from(sections).forEach((section, i) => {
			const rect = section.getBoundingClientRect();

			if (rect.top > threshold || rect.bottom < threshold) return;

			const activeAnchors = Array.from(links)
				.filter((link, ii) => link.classList.toggle(activeClass, i === ii));

			if (!cb) return;

			cb(activeAnchors[0]);
		});
	}

	scrollContainer.addEventListener('scroll', handleScroll);
}

/**
 * Spy section and set correspond link to active
 * @param {object}
 * |_@param {node} el - the DOM element to watch
 * |_@param {string} fixedClass
 * |_@param {float} threshold
 * |_@param {function} cb
 *   |_@param {bool} fixed
 */
export function spyNav({
	el,
	nav = el,
	fixedClass = 'is-fixed',
	threshold = 1,
	cb,
}) {
	const delta = 5;

	// State variables
	let enabled = false;
	let lastY = 0;

	// Event handler functions
	const handleScroll = debounce(() => {
		// Only peek header if nav to spy is fixed/enabled
		if (!enabled) {
			lastY = 0;

			return;
		}

		const thisY = window.pageYOffset;

		// Add delta for sensitivity threshold
		if (Math.abs(thisY - lastY) < delta) {
			return;
		}

		cb(thisY > lastY);

		lastY = thisY;
	}, 250);
	function handleFixObserver(entries) {
		// Toggle enabled state variable to conditionally process header peek on scroll up
		enabled = el.classList.toggle(fixedClass, entries[0].boundingClientRect.top < 0);

		cb(enabled);
	}

	// Add event listeners
	window.addEventListener('scroll', handleScroll);
	const fixObserver = new IntersectionObserver(handleFixObserver, { threshold });
	fixObserver.observe(nav);
}

export function spyImages() {
	const images = document.querySelectorAll('[loading="lazy"]');

	function handleObserver(entries, observer) {
		entries
			.filter(entry => entry.isIntersecting)
			.forEach(entry => {
				const image = entry.target;
				const { src } = image.dataset;

				switch (image.nodeName) {
				case 'DIV':
					image.style.backgroundImage = `url(${src})`;

					break;
				case 'SOURCE':
					image.srcset = src;

					break;
				default:
					image.src = src;
				}

				image.style.opacity = 1;
				image.removeAttribute('loading');
				observer.unobserve(image);
			});
	}

	const nativeSupport = 'loading' in HTMLImageElement.prototype;
	const observer = new IntersectionObserver(handleObserver, { threshold: 0.1 });

	images.forEach(image => {
		if (!nativeSupport || image.nodeName !== 'IMG') {
			observer.observe(image);
		} else {
			image.src = image.dataset.src;
            image.style.opacity = 1;
		}
	});
}
