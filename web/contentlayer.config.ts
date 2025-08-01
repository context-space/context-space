import { defineDocumentType, makeSource } from "contentlayer2/source-files"
import remarkGfm from "remark-gfm"

export const Blog = defineDocumentType(() => ({
  name: "Blog",
  filePathPattern: `**/*.md`,
  contentType: "mdx",
  fields: {
    title: { type: "string", required: true },
    description: { type: "string", required: true },
    publishedAt: { type: "date", required: true },
    category: { type: "string", required: false },
    author: { type: "string", required: true },
    image: { type: "string", required: true },
    featured: { type: "number", required: false, default: 0 },
  },
  computedFields: {
    url: { type: "string", resolve: post => `/blogs/${post._raw.flattenedPath}` },
    id: { type: "string", resolve: post => post._raw.flattenedPath },
  },
}))

export default makeSource({
  contentDirPath: "content/blogs",
  documentTypes: [Blog],
  mdx: {
    remarkPlugins: [remarkGfm],
  },
})
