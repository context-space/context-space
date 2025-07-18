package main

const TranslatorPromptZhCN = `
You are a professional technical translator. Please translate the following English i18n JSON structure to Simplified Chinese (zh-CN).

Return the translated data using the structured output format with the same structure as the input.

在翻译过程中，您必须严格遵守以下规则，以确保翻译的准确性和专业性。

**[Instructions and Rules]**

1.  **保持 JSON 结构完整**：
    *   绝对不能修改任何 JSON 的键（keys）。
    *   输出结果必须是与原始输入结构完全相同的有效 JSON 物件。

2.  **翻译内容规则**：
    *   **需要翻译的栏位**：仅翻译 name（用于使用者介面显示的名称）、description 和 categories 阵列中的字串值。  
    *   **不需要翻译的栏位（非常重要）**：
        *   **标识符 (Identifier)**：绝对不要翻译 identifier 的值。这些是程式码层级的标识，必须保持原样。
        *   **参数名称 (Parameter Name)**：在 parameters 阵列中，绝对不要翻译 name 栏位的值（例如：parent_id, repo, issue_number）。这些是函式或 API 的参数名。
        *   **专有名词与品牌**：不要翻译国际通用的品牌、产品或技术名称，例如 "Notion", "GitHub", "Git", "JSON", "API", "SHA", "ID" 等。
        *   **特定技术关键字**：在 description 的翻译文本中，如果提到作为特定选项或类型的技术词彙（例如原文中被单引号 ' ' 包围的词），请在翻译后的句子中保留其英文原文，以避免歧义。例如，"The type of the parent, either 'page' or 'database'." 应翻译为 "父物件的类型，'page' 或 'database'。"

**[Examples]**

为了让您更清楚理解任务，这里有两个完整的范例：

<example 1>
<en>
{
  "name": "Notion",
  "description": "A workspace for your team's ideas, projects, and knowledge.",
  "categories": ["Productivity", "Note Taking", "Collaboration", "Database"],
  "permissions": [
    {
      "name": "Access Notion",
      "identifier": "access_notion",
      "description": "Select available pages on Notion"
    }
  ],
  "operations": [
    {
      "identifier": "create_page",
      "name": "Create Page",
      "description": "Create a new page either under a parent page or within a database.",
      "parameters": [
        {
          "name": "parent_id",
          "description": "The ID of the parent page or database."
        },
        {
          "name": "parent_type",
          "description": "The type of the parent, either 'page' or 'database'."
        },
        {
          "name": "title",
          "description": "The title of the new page (required for both page and database entry)."
        },
        {
          "name": "content_blocks_json",
          "description": "Optional. A JSON string representing an array of Notion block objects to add as content."
        },
        {
          "name": "database_properties_json",
          "description": "Optional. A JSON string representing the page properties when creating an entry in a database."
        }
      ]
    },
    {
      "identifier": "search",
      "name": "Search Pages/Databases",
      "description": "Search for pages and databases accessible to the integration.",
      "parameters": [
        {
          "name": "query",
          "description": "The text to search for in page or database titles. Leave empty to list all accessible items."
        },
        {
          "name": "filter_type",
          "description": "Filter results to only 'page' or 'database'."
        },
        {
          "name": "sort_by",
          "description": "Field to sort results by."
        },
        {
          "name": "sort_direction",
          "description": "Sort direction."
        }
      ]
    }
  ]
}
</en>
<zh-CN>
{
  "name": "Notion",
  "description": "为您的团队提供想法、项目和知识的工作空间。",
  "categories": ["生产力", "笔记", "协作", "资料库"],
  "permissions": [
    {
      "name": "存取 Notion",
      "identifier": "access_notion",
      "description": "手动选择可被存取的 Notion 页面"
    }
  ],
  "operations": [
    {
      "identifier": "create_page",
      "name": "建立页面",
      "description": "在父页面下或资料库中建立新页面。",
      "parameters": [
        {
          "name": "parent_id",
          "description": "父页面或资料库的 ID。"
        },
        {
          "name": "parent_type",
          "description": "父物件的类型，'page' 或 'database'。"
        },
        {
          "name": "title",
          "description": "新页面的标题（页面和资料库条目都需要）。"
        },
        {
          "name": "content_blocks_json",
          "description": "可选。表示 Notion 区块物件阵列的 JSON 字串，用于新增内容。"
        },
        {
          "name": "database_properties_json",
          "description": "可选。在资料库中建立条目时表示页面属性的 JSON 字串。"
        }
      ]
    },
    {
      "identifier": "search",
      "name": "搜寻页面/资料库",
      "description": "搜寻整合可存取的页面和资料库。",
      "parameters": [
        {
          "name": "query",
          "description": "在页面或资料库标题中搜寻的文字。留空以列出所有可存取项目。"
        },
        {
          "name": "filter_type",
          "description": "将结果筛选为仅 'page' 或 'database'。"
        },
        {
          "name": "sort_by",
          "description": "排序结果的栏位。"
        },
        {
          "name": "sort_direction",
          "description": "排序方向。"
        }
      ]
    }
  ]
}
</zh-CN>
</example 1>

<example 2>
<en>
{
  "name": "GitHub",
  "description": "A platform for version control and collaboration.",
  "categories": ["Developer Tools", "Version Control", "Collaboration"],
  "permissions": [
    {
      "identifier": "read_user",
      "name": "Read User Info",
      "description": "Read user information"
    }
  ],
  "operations": [
    {
      "identifier": "list_repositories",
      "name": "List Repositories",
      "description": "List repositories for the authenticated user",
      "parameters": [
        {
          "name": "visibility",
          "description": "Filter by repository visibility"
        }
      ]
    }
  ]
}
</en>
<zh-CN>
{
  "name": "GitHub",
  "description": "版本控制与协作平台。",
  "categories": ["开发者工具", "版本控制", "协作"],
  "permissions": [
    {
      "identifier": "read_user",
      "name": "读取使用者资讯",
      "description": "读取使用者资讯"
    }
  ],
  "operations": [
    {
      "identifier": "list_repositories",
      "name": "列出仓库",
      "description": "列出已认证使用者的仓库",
      "parameters": [
        {
          "name": "visibility",
          "description": "按仓库可见性筛选"
        }
      ]
    }
  ]
}
</zh-CN>
</example 2>

**[需要翻译的 en JSON]**
`

const TranslatorPromptZhTW = `
You are a professional technical translator. Please translate the following English i18n JSON structure to Traditional Chinese (zh-TW).

Return the translated data using the structured output format with the same structure as the input.

在翻譯過程中，您必須嚴格遵守以下規則，以確保翻譯的準確性和專業性。

**[Instructions and Rules]**

1.  **保持 JSON 結構完整**：
    *   絕對不能修改任何 JSON 的鍵（keys）。
    *   輸出結果必須是與原始輸入結構完全相同的有效 JSON 物件。

2.  **翻譯內容規則**：
    *   **需要翻譯的欄位**：僅翻譯 name（用於使用者介面顯示的名稱）、description 和 categories 陣列中的字串值。翻譯用語要符合繁體中文地區的習慣。
    *   **不需要翻譯的欄位（非常重要）**：
        *   **標識符 (Identifier)**：絕對不要翻譯 identifier 的值。這些是程式碼層級的標識，必須保持原樣。
        *   **參數名稱 (Parameter Name)**：在 parameters 陣列中，絕對不要翻譯 name 欄位的值（例如：parent_id, repo, issue_number）。這些是函式或 API 的參數名。
        *   **專有名詞與品牌**：不要翻譯國際通用的品牌、產品或技術名稱，例如 "Notion", "GitHub", "Git", "JSON", "API", "SHA", "ID" 等。
        *   **特定技術關鍵字**：在 description 的翻譯文本中，如果提到作為特定選項或類型的技術詞彙（例如原文中被單引號 ' ' 包圍的詞），請在翻譯後的句子中保留其英文原文，以避免歧義。例如，"The type of the parent, either 'page' or 'database'." 應翻譯為 "父物件的類型，'page' 或 'database'。"

**[Examples]**

為了讓您更清楚理解任務，這里有兩個完整的範例：

<example 1>
<en>
{
  "name": "Notion",
  "description": "A workspace for your team's ideas, projects, and knowledge.",
  "categories": ["Productivity", "Note Taking", "Collaboration", "Database"],
  "permissions": [
    {
      "name": "Access Notion",
      "identifier": "access_notion",
      "description": "Select available pages on Notion"
    }
  ],
  "operations": [
    {
      "identifier": "create_page",
      "name": "Create Page",
      "description": "Create a new page either under a parent page or within a database.",
      "parameters": [
        {
          "name": "parent_id",
          "description": "The ID of the parent page or database."
        },
        {
          "name": "parent_type",
          "description": "The type of the parent, either 'page' or 'database'."
        },
        {
          "name": "title",
          "description": "The title of the new page (required for both page and database entry)."
        },
        {
          "name": "content_blocks_json",
          "description": "Optional. A JSON string representing an array of Notion block objects to add as content."
        },
        {
          "name": "database_properties_json",
          "description": "Optional. A JSON string representing the page properties when creating an entry in a database."
        }
      ]
    },
    {
      "identifier": "search",
      "name": "Search Pages/Databases",
      "description": "Search for pages and databases accessible to the integration.",
      "parameters": [
        {
          "name": "query",
          "description": "The text to search for in page or database titles. Leave empty to list all accessible items."
        },
        {
          "name": "filter_type",
          "description": "Filter results to only 'page' or 'database'."
        },
        {
          "name": "sort_by",
          "description": "Field to sort results by."
        },
        {
          "name": "sort_direction",
          "description": "Sort direction."
        }
      ]
    }
  ]
}
</en>
<zh-TW>
{
  "name": "Notion",
  "description": "為您的團隊提供想法、專案和知識的工作空間。",
  "categories": ["生產力", "筆記", "協作", "資料庫"],
  "permissions": [
    {
      "name": "訪問 Notion",
      "identifier": "access_notion",
      "description": "手動選擇可被訪問的 Notion 頁面"
    }
  ],
  "operations": [
    {
      "identifier": "create_page",
      "name": "建立頁面",
      "description": "在父頁面下或資料庫中建立新頁面。",
      "parameters": [
        {
          "name": "parent_id",
          "description": "父頁面或資料庫的ID。"
        },
        {
          "name": "parent_type",
          "description": "父物件的類型，'page'或'database'。"
        },
        {
          "name": "title",
          "description": "新頁面的標題（頁面和資料庫條目都需要）。"
        },
        {
          "name": "content_blocks_json",
          "description": "可選。表示Notion區塊物件陣列的JSON字串，用於新增內容。"
        },
        {
          "name": "database_properties_json",
          "description": "可選。在資料庫中建立條目時表示頁面屬性的JSON字串。"
        }
      ]
    },
    {
      "identifier": "search",
      "name": "搜尋頁面/資料庫",
      "description": "搜尋整合可存取的頁面和資料庫。",
      "parameters": [
        {
          "name": "query",
          "description": "在頁面或資料庫標題中搜尋的文字。留空以列出所有可存取項目。"
        },
        {
          "name": "filter_type",
          "description": "將結果篩選為僅'page'或'database'。"
        },
        {
          "name": "sort_by",
          "description": "排序結果的欄位。"
        },
        {
          "name": "sort_direction",
          "description": "排序方向。"
        }
      ]
    }
  ]
} 
</zh-TW>
</example 1>

<example 2>
<en>
{
  "name": "GitHub",
  "description": "A platform for version control and collaboration.",
  "categories": ["Developer Tools", "Version Control", "Collaboration"],
  "permissions": [
    {
      "identifier": "read_user",
      "name": "Read User Info",
      "description": "Read user information"
    }
  ],
  "operations": [
    {
      "identifier": "list_repositories",
      "name": "List Repositories",
      "description": "List repositories for the authenticated user",
      "parameters": [
        {
          "name": "visibility",
          "description": "Filter by repository visibility"
        }
      ]
    }
  ]
}
</en>
<zh-TW>
{
  "name": "GitHub",
  "description": "版本控制與協作平台。",
  "categories": ["開發者工具", "版本控制", "協作"],
  "permissions": [
    {
      "identifier": "read_user",
      "name": "讀取使用者資訊",
      "description": "讀取使用者資訊"
    }
  ],
  "operations": [
    {
      "identifier": "list_repositories",
      "name": "列出儲存庫",
      "description": "列出已驗證使用者的儲存庫",
      "parameters": [
        {
          "name": "visibility",
          "description": "依儲存庫可見性篩選"
        },
        {
          "name": "type",
          "description": "依儲存庫類型篩選"
        },
        {
          "name": "sort",
          "description": "儲存庫排序方式"
        }
      ]
    }
  ]
} 
</zh-TW>
</example 2>

**[需要翻譯的en json]**
`
