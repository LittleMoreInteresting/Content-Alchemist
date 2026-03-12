# Content Alchemist 版本发布脚本 (Windows PowerShell)
# 使用方法: .\scripts\release.ps1 [version_type]
# version_type: patch (默认), minor, major, alpha, beta

param(
    [string]$VersionType = "patch"
)

$ErrorActionPreference = "Stop"

# 颜色定义
$Red = "Red"
$Green = "Green"
$Yellow = "Yellow"
$Blue = "Cyan"

$CurrentBranch = git branch --show-current

# 检查是否在 main 或 master 分支
if ($CurrentBranch -ne "main" -and $CurrentBranch -ne "master") {
    Write-Host "❌ 错误: 必须在 main 或 master 分支上发布版本" -ForegroundColor $Red
    Write-Host "当前分支: $CurrentBranch"
    exit 1
}

# 检查工作目录是否干净
$status = git status --porcelain
if ($status) {
    Write-Host "❌ 错误: 工作目录不干净，请先提交或暂存更改" -ForegroundColor $Red
    git status --short
    exit 1
}

# 获取当前版本
try {
    $CurrentVersion = git describe --tags --abbrev=0
} catch {
    $CurrentVersion = "v0.0.0"
}
Write-Host "当前版本: $CurrentVersion" -ForegroundColor $Blue

# 解析版本号
$CurrentVersion = $CurrentVersion -replace '^v', ''
$versionParts = $CurrentVersion -split '\.'
$Major = [int]$versionParts[0]
$Minor = [int]$versionParts[1]
$Patch = [int]($versionParts[2] -split '-')[0]

# 计算新版本
switch ($VersionType) {
    "major" { $NewVersion = "$($Major + 1).0.0" }
    "minor" { $NewVersion = "$Major.$($Minor + 1).0" }
    "patch" { $NewVersion = "$Major.$Minor.$($Patch + 1)" }
    "alpha" { $NewVersion = "$Major.$Minor.$Patch-alpha.1" }
    "beta" { $NewVersion = "$Major.$Minor.$Patch-beta.1" }
    default {
        # 如果输入是完整的版本号
        if ($VersionType -match '^v?\d+\.\d+\.\d+') {
            $NewVersion = $VersionType -replace '^v', ''
        } else {
            Write-Host "❌ 错误: 无效的版本类型 '$VersionType'" -ForegroundColor $Red
            Write-Host "支持的类型: major, minor, patch, alpha, beta, 或完整版本号 (如 1.2.3)"
            exit 1
        }
    }
}

$NewTag = "v$NewVersion"

# 确认
Write-Host ""
Write-Host "即将创建新版本: $NewTag" -ForegroundColor Yellow
Write-Host ""
$Confirm = Read-Host "确认发布? (y/N)"

if ($Confirm -ne "y" -and $Confirm -ne "Y") {
    Write-Host "已取消发布" -ForegroundColor Yellow
    exit 0
}

Write-Host ""
Write-Host "🚀 开始发布流程..." -ForegroundColor Blue

# 更新 wails.json 中的版本
if (Test-Path "wails.json") {
    Write-Host "📦 更新 wails.json 版本..." -ForegroundColor Blue
    $content = Get-Content wails.json -Raw
    $content = $content -replace '"productVersion": ".*?"', "`"productVersion`": `"$NewVersion`""
    Set-Content wails.json $content
}

# 更新 frontend/package.json 中的版本
if (Test-Path "frontend/package.json") {
    Write-Host "📦 更新 package.json 版本..." -ForegroundColor Blue
    $content = Get-Content frontend/package.json -Raw
    $content = $content -replace '"version": ".*?"', "`"version`": `"$NewVersion`""
    Set-Content frontend/package.json $content
}

# 提交版本更新（如果有变更）
$status = git status --porcelain
if ($status) {
    Write-Host "📝 提交版本更新..." -ForegroundColor Blue
    git add wails.json frontend/package.json
    git commit -m "chore(release): bump version to $NewTag"
    git push origin $CurrentBranch
}

# 创建并推送标签
Write-Host "🏷️ 创建标签 $NewTag..." -ForegroundColor Blue
git tag -a $NewTag -m "Release $NewTag"
git push origin $NewTag

Write-Host ""
Write-Host "✅ 版本 $NewTag 发布成功!" -ForegroundColor Green
Write-Host ""
Write-Host "🔄 GitHub Actions 正在构建发布..." -ForegroundColor Blue

# 获取仓库 URL
$remoteUrl = git remote get-url origin
$repoPath = $remoteUrl -replace '.*github.com[:/]', '' -replace '\.git$', ''
Write-Host "查看进度: https://github.com/$repoPath/actions"
Write-Host ""
Write-Host "注意:" -ForegroundColor Yellow
Write-Host "  - Release 初始为 Draft 状态，构建完成后会自动发布"
Write-Host "  - 你可以在 GitHub 上查看并编辑 Release 说明"
