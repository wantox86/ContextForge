export interface Rule {
  ID: number
  TemplateID: number
  Category: 'coding_style' | 'testing' | 'security' | 'naming' | 'architecture'
  Title: string
  Content: string
  Priority: number
  Enabled: boolean
}

export interface Template {
  ID: number
  Name: string
  Description: string
  Stack: string[]
  Rules: Rule[]
  IsBuiltIn: boolean
  CreatedAt: string
  UpdatedAt: string
}

export interface Project {
  ID: number
  Name: string
  Path: string
  DetectedStack: string[]
  AppliedTemplates: number[]
  CreatedAt: string
}

export interface DetectionResult {
  stack: string[]
  confidence: number
  indicators: string[]
}

export interface TokenWarning {
  Model: string
  Tokens: number
  WarnAt: number
  Max: number
  Exceeds: boolean
}
