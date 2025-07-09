import { redirect } from "next/navigation"

export default function DocsPage() {
  return redirect("https://github.com/context-space/context-space/blob/main/README.md")
}
