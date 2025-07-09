import { atom } from "jotai"

// Atom to store the selected operation name that should be filled into the playground input
export const selectedOperationAtom = atom<string>("")

// Atom to trigger input update in playground
export const triggerInputUpdateAtom = atom<number>(0)
