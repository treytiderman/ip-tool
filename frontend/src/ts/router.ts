import { writable } from 'svelte/store'

export {
    pages,
    pageStore,
    setPage,
    setHomePage,
}

let pages = [
    { name: "IPv4 Presets" },
    { name: "Create Preset" },
    { name: "Edit Preset" },
    { name: "Edit Interface" },
    { name: "Settings" },
]

let pageStore = writable(pages[0])

function setHomePage() {
    setPage("IPv4 Presets")
}

function setPage(name: string) {
    const page = pages.find(p => p.name === name)
    if (page) {
        pageStore.set(page)
    } else {
        console.warn(`Page with name "${name}" not found.`)
    }
}
