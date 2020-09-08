import { item } from "./item.js";
import { currentlySelectedIdx } from "./store.js";

export function getLocationLastElement() {
  var loc = window.location.pathname;
  var parts = loc.split("/");
  var lastIdx = parts.length - 1;
  return parts[lastIdx];
}

export function getLocationLastElementWithHash() {
  var loc = window.location.pathname;
  var parts = loc.split("/");
  var lastIdx = parts.length - 1;
  return parts[lastIdx] + window.location.hash;
}

// TODO: maybe move to item.js
// remembers which toc items are expanded, by their index
let tocItemIdxExpanded = [];

export function isTocItemExpanded(idx) {
  for (let i of tocItemIdxExpanded) {
    if (i === idx) {
      return true;
    }
  }
  return false;
}

function setIsExpandedUpwards(idx) {
  const tocItem = gTocItems[idx];
  tocItemIdxExpanded.push(idx);
  // console.log(`idx: ${idx}, title: ${tocItem[4]}`)
  const newIdx = item.parentIdx(tocItem);
  if (newIdx != -1) {
    setIsExpandedUpwards(newIdx);
  }
}

export function findTocIdxForCurrentURL() {
  const currURI = getLocationLastElementWithHash();
  const n = gTocItems.length;
  let tocItem, uri;
  for (let idx = 0; idx < n; idx++) {
    tocItem = gTocItems[idx];
    uri = item.url(tocItem);
    if (uri === currURI) {
      return idx;
    }
  }
  return -1;
}

export function setTocExpandedForCurrentURL() {
  tocItemIdxExpanded = [];
  const idx = findTocIdxForCurrentURL();
  if (idx == -1) {
    return 0;
  }

  currentlySelectedIdx.set(idx);
  setIsExpandedUpwards(idx);
  return idx;
}

// returns a debouncer function. Usage:
// var debouncer = makeDebouncer(250);
// function fn() { ... }
// debouncer(fn)
export function makeDebouncer(timeInMs) {
  let interval;
  return function (f) {
    clearTimeout(interval);
    interval = setTimeout(() => {
      interval = null;
      f();
    }, timeInMs);
  };
}

// d is json data to be logged to logflare
const logflareURL = "https://api.logflare.app/logs?api_key=EMVM15g1Nf_M&source=b72adcc3-d077-46cc-92c3-2f0c55cf3b69";

export function logflare(d) {
  const data = {
    headers: {
      'Content-Type': 'application/json; charset=utf-8'
    },
    method: "post",
    body: JSON.stringify(d),
  };
  function thenFn(res) {
    console.log(`sent log to logflare`);
    console.log(res);
  }
  function catchFn(res) {
    console.log(`exception in logflare`);
  }
  fetch(logflareURL, data)
  .then(thenFn)
  .catch(catchFn);
}
