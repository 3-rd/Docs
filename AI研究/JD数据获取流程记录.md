# JD 数据获取流程记录

> 记录各公司招聘平台 API 的访问方式、参数结构、验证方式
> 更新时间：2026-04-08

---

## 1. 腾讯（careers.tencent.com）

### 平台说明
腾讯社招平台，有两个主要域：
- `careers.tencent.com` — 主招聘页
- `tencentcareer.qq.com` — 历史域名（部分重定向）

### API 端点
```
GET https://careers.tencent.com/tencentcareer/api/post/Query
```

### 参数说明
| 参数 | 说明 | 示例 |
|------|------|------|
| `timestamp` | 当前毫秒时间戳（绕过缓存） | `1775657290833` |
| `categoryId` | 岗位类别 ID，AI 相关固定值 | `40001005` |
| `pageIndex` | 页码（注意是 pageIndex 不是 pageNum） | `1` |
| `pageSize` | 每页数量 | `10` |
| `language` | 语言 | `zh-cn` |
| `area` | 地区 | `cn` |
| `attrId` | 属性筛选 | `1`（固定） |

### 请求示例
```
GET https://careers.tencent.com/tencentcareer/api/post/Query?timestamp=1775657290833&categoryId=40001005&pageIndex=1&pageSize=10&language=zh-cn&area=cn&attrId=1
```

### 验证方式
- 无需 Cookie 或 Token
- 直接访问即可，可能触发 IP 限流（返回 500）
- 限流后只能等几个小时后重试
- ByPostId 详情接口：`GET /tencentcareer/api/post/ByPostId?postId={postId}&language=zh-cn`

### 数据结构
```json
{
  "Code": 200,
  "Data": {
    "Count": 272,       // 总数
    "Posts": [           // 注意是 Posts 大写
      {
        "PostId": "xxx",
        "RecruitPostName": "岗位名",
        "LocationName": "城市",
        "BGName": "BG",
        "RequireWorkYearsName": "三年以上",
        "Responsibility": "职责（列表页无完整正文）",
        "Requirement": "要求（列表页无完整正文）"
      }
    ]
  }
}
```

**重要**：列表页 `Responsibility` 和 `Requirement` 字段内容不完整，完整 JD 需要通过 `ByPostId` 逐个调详情接口获取。

### 详情接口（ByPostId）
```
GET https://careers.tencent.com/tencentcareer/api/post/ByPostId?postId={postId}&language=zh-cn
```
返回完整 `Requirement` + `Responsibility` 字段。

---

## 2. 阿里（talent-holding.alibaba.com）

### 平台说明
阿里巴巴内部招聘平台（内网风格），需要登录态 Cookie。

### API 端点
```
POST https://talent-holding.alibaba.com/position/search?_csrf={csrf_token}
```

### 参数说明
| 参数 | 说明 | 示例 |
|------|------|------|
| `_csrf` | CSRF Token，需从页面获取 | `894c8ed9-6c28-4ce0-9dbc-b8f93bb2780a` |
| `keyword` | 搜索关键词 | `大模型` |
| `pageSize` | 每页数量 | `10` |
| `currentPage` | 页码 | `1` |

### 请求头
```
Content-Type: application/json
Cookie: SESSION=xxx; XSRF-TOKEN=xxx; ...
X-CSRF-Token: {csrf_token}
Referer: https://talent-holding.alibaba.com/position/search?_csrf=xxx
```

### 请求体（JSON）
```json
{"keyword": "大模型", "pageSize": 10, "currentPage": 1}
```

### Cookie 获取方式
在浏览器登录后，F12 → Application → Cookies → 复制 `SESSION` 和 `XSRF-TOKEN`。

### 数据结构
```json
{
  "content": {
    "totalCount": 73,
    "datas": [
      {
        "positionId": "xxx",
        "name": "岗位名",
        "workLocations": ["杭州"],
        "experience": {"from": 3, "to": null},
        "degree": "bachelor",
        "requirement": "完整要求正文",
        "description": "完整职责正文"
      }
    ]
  }
}
```

---

## 3. 阿里云（careers.aliyun.com）

### 平台说明
阿里云对外招聘平台，公开访问。

### API 端点
```
POST https://careers.aliyun.com/position/search?_csrf={csrf_token}
```

### 参数说明（POST Body）
```json
{
  "channel": "group_official_site",
  "language": "zh",
  "batchId": "",
  "categories": "130",
  "deptCodes": [],
  "key": "",
  "pageIndex": 1,
  "pageSize": 10,
  "regions": "",
  "subCategories": "136"
}
```

| 字段 | 说明 |
|------|------|
| `categories` | 岗位大类，AI 相关固定 `130` |
| `subCategories` | 子类，`130,747,100000053,136,409,408` |
| `pageIndex` | 页码 |
| `pageSize` | 每页数量 |

### 请求头
```
Content-Type: application/json
Cookie: 从浏览器复制完整 cookie
X-CSRF-Token: 从 URL 获取
Referer: https://careers.aliyun.com/off-campus/position-list?lang=zh
```

### Cookie 获取方式
在浏览器打开 `careers.aliyun.com` 并登录，F12 → Network → 任意请求 → 复制完整 Cookie 头。

### 数据结构
```json
{
  "success": true,
  "content": {
    "totalCount": 276,
    "datas": [
      {
        "id": 100004043002,
        "trackId": "SSPxxx",
        "name": "岗位名",
        "workLocations": ["北京", "杭州"],
        "requirement": "完整要求正文",
        "description": "完整职责正文",
        "experience": {"from": 3},
        "degree": "master"
      }
    ]
  }
}
```

### 详情页 URL 格式
```
https://careers.aliyun.com/off-campus/position-detail?positionId={ID}&track_id={trackId}
```

---

## 4. 蚂蚁集团（talent.antgroup.com）

### 平台说明
蚂蚁集团招聘平台，需要登录态。

### API 端点
```
POST https://hrcareersweb.antgroup.com/api/social/position/search?ctoken={ctoken}
```

### 参数说明（POST Body）
```json
{
  "regions": "",
  "categories": "130",
  "subCategories": "130,747,100000053,136,409,408",
  "bgCode": "",
  "socialQrCode": "",
  "pageIndex": 1,
  "pageSize": 10,
  "channel": "group_official_site",
  "language": "zh"
}
```

| 字段 | 说明 |
|------|------|
| `categories` | 固定 `130` |
| `subCategories` | AI 相关子类组合 |
| `pageIndex` | 页码 |
| `pageSize` | 每页数量 |

### 请求头
```
Content-Type: application/json;charset=UTF-8
Cookie: 从浏览器复制（包含 SESSION、ctoken 等）
Referer: https://talent.antgroup.com/
Origin: https://talent.antgroup.com
```

### Cookie 获取方式
在浏览器打开 `talent.antgroup.com` 并登录，F12 → Network → 任意 POST 请求 → 复制完整 Request Headers 中的 Cookie。

### 数据结构（注意：content 是数组不是对象）
```json
{
  "success": true,
  "content": [
    {
      "id": 26022708835678,
      "name": "蚂蚁国际-AI工程师-全球技术",
      "workLocations": ["上海", "杭州", "深圳"],
      "requirement": "完整要求正文",
      "description": "完整职责正文",
      "experience": {"from": 3, "to": null},
      "degree": "bachelor",
      "department": "蚂蚁国际",
      "featureTagList": ["AI智能体", "Java", "LLM推理服务", ...]
    }
  ],
  "totalCount": 285
}
```

### 详情页 URL 格式
```
https://talent.antgroup.com/newcloud/position-detail?positionId={id}
```

---

## 5. 字节跳动（jobs.bytedance.com）

### 平台说明
字节跳动社招平台，有 SigV3 签名机制，比普通 Cookie 方案更严格。

### API 端点
```
POST https://jobs.bytedance.com/api/v1/search/job/posts?_signature={signature}
```

### 参数说明

**URL Query Parameters：**
| 参数 | 说明 | 示例 |
|------|------|------|
| `_signature` | **签名，时效极短**（几分钟内），每次刷新页面会变 | `xWTqQwAAAABJgZfT5ss0B8Vk6lAAKy8` |
| `keyword` | 搜索关键词 | `大模型`（留空则搜索全部） |
| `limit` | 每页数量 | `50`（建议用50减少请求数） |
| `offset` | 偏移量 | `0`、`50`、`100`... |
| `job_category_id_list` | 岗位类别 ID，AI 相关固定值 | `6704215956018694411,6704215862557018372,6704219534724696331` |
| `portal_type` | 固定值 | `2` |
| `portal_entrance` | 固定值 | `1` |

**POST Body（JSON）：**
```json
{
  "keyword": "",
  "limit": 50,
  "offset": 0,
  "job_category_id_list": [
    "6704215956018694411",
    "6704215862557018372",
    "6704219534724696331"
  ],
  "tag_id_list": [],
  "location_code_list": [],
  "subject_id_list": [],
  "recruitment_id_list": [],
  "portal_type": 2,
  "job_function_id_list": [],
  "storefront_id_list": [],
  "portal_entrance": 1
}
```

### 请求头（必须完整）
```
Content-Type: application/json
Cookie: locale=zh-CN; channel=office; platform=pc; s_v_web_id=xxx; device-id=xxx
Referer: https://jobs.bytedance.com/experienced/position
Origin: https://jobs.bytedance.com
Portal-Channel: office
Portal-Platform: pc
Accept: application/json, text/plain, */*
```

### Cookie 获取方式
1. 在浏览器打开 `https://jobs.bytedance.com` 并登录
2. F12 → Network → 刷新页面或搜索 → 找到 `job/posts` 请求
3. 复制 Request Headers 中的完整 `Cookie` 字段

### _signature 获取方式
1. 在上述请求中复制 `_signature` 参数的值
2. **注意**：`signature` 时效极短（几分钟），建议一次性快速抓完所有页
3. 如果中途失败，需要刷新页面重新获取新的 `signature`

### 数据结构
```json
{
  "code": 0,
  "data": {
    "count": 2138,        // 总数（固定返回）
    "job_post_list": [
      {
        "id": "7621153447852755253",
        "title": "支付风控算法工程师-国际支付",
        "description": "完整职责正文（带换行）",
        "requirement": "完整要求正文（带换行）",
        "city_info": {
          "code": "CT_128",
          "name": "深圳"
        },
        "city_list": [
          {"code": "CT_128", "name": "深圳"}
        ],
        "job_category": {
          "id": "6704215956018694411",
          "name": "算法"
        },
        "publish_time": 1774438151591
      }
    ]
  }
}
```

### 注意事项
- **`_signature` 时效短**：建议用 `limit=50` 一次多抓几页，减少请求次数，缩短总耗时
- **经验要求解析困难**：字节 JD 的经验要求多在正文中表述（而非固定字段），精确筛选需读正文
- **总数在 count 字段**：可通过第一页响应中的 `count` 字段获取总岗位数

### 详情页 URL 格式
```
https://jobs.bytedance.com/experienced/position/detail/{job_id}
```

---

## 通用经验总结

### 快速获取 Cookie 的方法
1. 在浏览器登录对应平台
2. F12 打开开发者工具 → Network 标签
3. 刷新页面，随便点一个请求
4. 找到 Request Headers → 复制完整的 `Cookie` 字段
5. 对于 CSRF Token，通常在 URL query string 或 Cookie 里都有

### 反爬/限流处理
- 腾讯：`/Query` 接口最敏感，限流后返回 500，等几小时重试
- 阿里/阿里云/蚂蚁：通常 Cookie 有效期内稳定可用

### 数据完整性
| 平台 | 列表页字段 | 详情接口 | 备注 |
|------|-----------|---------|------|
| 腾讯 | 少量 | ByPostId 单独调 | 列表页 Requirement 不完整 |
| 阿里 | 完整 | 不需要 | 直接可用 |
| 阿里云 | 完整 | 不需要 | 直接可用 |
| 蚂蚁 | 完整 | 不需要 | 直接可用 |
| 字节 | 完整 | 不需要 | 直接可用，经验要求在正文中 |

---

*最后更新：2026-04-08（新增字节跳动）*
