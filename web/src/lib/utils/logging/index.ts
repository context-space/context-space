import { createConsola } from "consola"

// 检查是否在浏览器环境
const isBrowser = typeof window !== "undefined"

// 检查是否在生产环境
const isProduction = process.env.NODE_ENV === "production"

// Server端logger - 始终启用
export const serverLogger = createConsola({
  level: 4,
  formatOptions: {
    colors: true,
    date: true,
    time: true,
    columns: 30,
  },
})

// Client端logger - 在生产环境下可以禁用
export const clientLogger = createConsola({
  level: isProduction ? 0 : 4, // 生产环境下禁用client日志
  formatOptions: {
    colors: true,
    date: true,
    time: true,
    columns: 30,
  },
})

// 根据环境自动选择logger
export const logger = isBrowser ? clientLogger : serverLogger

// 为了向后兼容，保留原有的导出
export { logger as default }
