import './polyfills';

// Utilities
import { spyImages } from './utilities/spy';

// Event handler functions
function handleDOMConentLoaded() {
    spyImages();
}

// Add event listeners
document.addEventListener('DOMContentLoaded', handleDOMConentLoaded);
