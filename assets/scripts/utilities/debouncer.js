
// be triggered. The function will be called after it stops being called for
// N milliseconds. If `immediate` is passed, trigger the function on the
// leading edge, instead of the trailing.

/* eslint-disable */
export default function (func, wait = 0, immediate) {
    let timeout;
    return function trottle() {
        const context = this;
        const args = arguments;
        const callNow = immediate && !timeout;
        function later() {
            timeout = null;
            if (!immediate) func.apply(context, args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
        if (callNow) func.apply(context, args);
    };
}
/* eslint-enable */
