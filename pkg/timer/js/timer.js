const id = (() => {
    let currentId = 0;
    const map = new WeakMap();
    return (object) => {
        if (!map.has(object)) {
            map.set(object, ++currentId);
        }
        return map.get(object);
    };
})();

globalThis.setTimeout = (callback, delay, ...arugments) => {
    var timerFunctionName = "__timerFunc__" + id(callback)
    globalThis[timerFunctionName] =  () => {
        callback(...arugments);
    }
    return globalThis.__setTimeout(timerFunctionName, delay)
};

globalThis.setInterval = (callback, delay, ...arugments) => {
    const timerFunctionName = "__timerFunc__" + id(callback)
    globalThis[timerFunctionName] =  () => {
        callback(...arugments);
    }
    return globalThis.__setInterval(timerFunctionName, delay)
}