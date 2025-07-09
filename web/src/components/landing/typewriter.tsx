import { useTranslations } from "next-intl"
import { useEffect, useState } from "react"

export function TypewriterText() {
  const t = useTranslations()
  const words = [t("hero.typewriter.word1"), t("hero.typewriter.word2"), t("hero.typewriter.word3")]

  const [currentWordIndex, setCurrentWordIndex] = useState(0)
  const [displayText, setDisplayText] = useState("")
  const [isDeleting, setIsDeleting] = useState(false)

  useEffect(() => {
    const word = words[currentWordIndex]
    const timeout = setTimeout(() => {
      if (!isDeleting) {
        // Typing
        if (displayText.length < word.length) {
          setDisplayText(word.slice(0, displayText.length + 1))
        } else {
          // Pause before deleting
          setTimeout(() => setIsDeleting(true), 1500)
        }
      } else {
        // Deleting
        if (displayText.length > 0) {
          setDisplayText(displayText.slice(0, -1))
        } else {
          // Move to next word
          setIsDeleting(false)
          setCurrentWordIndex(prev => (prev + 1) % words.length)
        }
      }
    }, isDeleting ? 50 : 100)

    return () => clearTimeout(timeout)
  }, [displayText, isDeleting, currentWordIndex, words])

  return (
    <p className="text-xl md:text-2xl lg:text-3xl text-neutral-600 dark:text-gray-300 leading-relaxed tracking-wide">
      {t("hero.typewriter.prefix")}
      {" "}
      <span className="text-primary font-medium">{displayText}</span>
      <span className="animate-pulse">|</span>
    </p>
  )
}
