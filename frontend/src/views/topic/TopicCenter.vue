<template>
  <div class="topic-center">
    <!-- 顶部标题栏 -->
    <div class="header">
      <div class="title-section">
        <h2>选题中心</h2>
        <p class="subtitle">发现热点，AI生成优质选题</p>
      </div>
      <div class="actions">
        <el-button type="primary" @click="refreshHotTrends" :loading="loadingHot">
          <el-icon><Refresh /></el-icon>
          刷新热点
        </el-button>
        <el-button type="success" @click="generateTopics" :loading="generating">
          <el-icon><MagicStick /></el-icon>
          AI生成选题
        </el-button>
      </div>
    </div>

    <!-- 标签页 -->
    <el-tabs v-model="activeTab" class="topic-tabs">
      <!-- AI推荐选题 -->
      <el-tab-pane label="AI推荐" name="ai">
        <div class="ai-recommend" v-loading="generating">
          <div v-if="aiTopics.length === 0 && !generating" class="empty-state">
            <el-empty description="点击右上角AI生成选题按钮开始">
              <template #image>
                <el-icon :size="60" color="#909399"><MagicStick /></el-icon>
              </template>
            </el-empty>
          </div>
          <div v-else class="topic-grid">
            <el-card
              v-for="topic in aiTopics"
              :key="topic.id"
              class="topic-card"
              :class="{ 'high-score': topic.score >= 80 }"
            >
              <div class="topic-header">
                <el-tag v-if="topic.score >= 80" type="danger" effect="dark">强烈推荐</el-tag>
                <el-tag v-else-if="topic.score >= 70" type="warning">推荐</el-tag>
                <span class="score">{{ topic.score?.toFixed(1) }}分</span>
              </div>
              
              <h4 class="topic-title">{{ topic.title }}</h4>
              
              <p class="topic-reason" v-if="topic.reason">{{ topic.reason }}</p>
              
              <div class="topic-angles" v-if="topic.angles?.length">
                <h5>切入角度：</h5>
                <div class="angle-tags">
                  <el-tag v-for="angle in topic.angles" :key="angle" size="small" effect="plain">
                    {{ angle }}
                  </el-tag>
                </div>
              </div>
              
              <div class="topic-keywords" v-if="topic.keywords?.length">
                <el-tag
                  v-for="kw in topic.keywords"
                  :key="kw"
                  type="info"
                  size="small"
                  effect="plain"
                >
                  {{ kw }}
                </el-tag>
              </div>
              
              <div class="topic-actions">
                <el-button type="primary" size="small" @click="createArticle(topic)">
                  立即创作
                </el-button>
                <el-button size="small" @click="saveToLibrary(topic)">
                  加入选题库
                </el-button>
              </div>
            </el-card>
          </div>
        </div>
      </el-tab-pane>

      <!-- 热点趋势 -->
      <el-tab-pane label="热点趋势" name="hot">
        <div class="hot-trends">
          <div class="platform-filters">
            <el-checkbox-group v-model="selectedPlatforms" @change="fetchHotTrends">
              <el-checkbox-button label="weibo">微博</el-checkbox-button>
              <el-checkbox-button label="zhihu">知乎</el-checkbox-button>
              <el-checkbox-button label="baidu">百度</el-checkbox-button>
              <el-checkbox-button label="toutiao">头条</el-checkbox-button>
            </el-checkbox-group>
          </div>
          
          <div class="trend-list" v-loading="loadingHot">
            <div v-if="hotTrends.length === 0" class="empty-state">
              <el-empty description="暂无热点数据" />
            </div>
            <div v-else class="platform-groups">
              <div
                v-for="(trends, platform) in groupedHotTrends"
                :key="platform"
                class="platform-section"
              >
                <h4 class="platform-title">
                  {{ platformNames[platform] || platform }}
                  <span class="count">({{ trends.length }})</span>
                </h4>
                <el-scrollbar height="400px">
                  <div class="trend-items">
                    <div
                      v-for="item in trends"
                      :key="item.id"
                      class="trend-item"
                      @click="openTrendUrl(item.url)"
                    >
                      <span class="rank" :class="{ 'top3': item.hotRank <= 3 }">
                        {{ item.hotRank }}
                      </span>
                      <span class="title">{{ item.title }}</span>
                      <span class="hot-value" v-if="item.hotValue">
                        {{ formatHotValue(item.hotValue) }}
                      </span>
                      <el-tag size="small" type="info" v-if="item.category">
                        {{ item.category }}
                      </el-tag>
                    </div>
                  </div>
                </el-scrollbar>
              </div>
            </div>
          </div>
        </div>
      </el-tab-pane>

      <!-- 我的选题库 -->
      <el-tab-pane label="我的选题库" name="library">
        <div class="topic-library">
          <div class="filter-bar">
            <el-radio-group v-model="libraryFilter" @change="loadLibrary">
              <el-radio-button label="">全部</el-radio-button>
              <el-radio-button label="pending">待处理</el-radio-button>
              <el-radio-button label="approved">已批准</el-radio-button>
              <el-radio-button label="processing">创作中</el-radio-button>
              <el-radio-button label="published">已发布</el-radio-button>
            </el-radio-group>
            <el-input
              v-model="searchKeyword"
              placeholder="搜索选题"
              :prefix-icon="Search"
              clearable
              @input="searchTopics"
              style="width: 200px"
            />
          </div>
          
          <el-table :data="libraryTopics" style="width: 100%" v-loading="loadingLibrary">
            <el-table-column prop="title" label="选题标题" min-width="250">
              <template #default="{ row }">
                <div class="topic-title-cell">
                  <span class="title">{{ row.title }}</span>
                  <div class="tags" v-if="row.keywords?.length">
                    <el-tag
                      v-for="kw in row.keywords.slice(0, 3)"
                      :key="kw"
                      size="small"
                      type="info"
                      effect="plain"
                    >
                      {{ kw }}
                    </el-tag>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="score" label="评分" width="100">
              <template #default="{ row }">
                <el-tag :type="getScoreType(row.score)">
                  {{ row.score?.toFixed(1) || '-' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="source" label="来源" width="100">
              <template #default="{ row }">
                <el-tag size="small" effect="plain">
                  {{ sourceNames[row.source] || row.source }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)">
                  {{ statusNames[row.status] || row.status }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="createdAt" label="创建时间" width="180">
              <template #default="{ row }">
                {{ formatDate(row.createdAt) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="{ row }">
                <el-button
                  v-if="row.status === 'pending'"
                  type="primary"
                  size="small"
                  @click="approveTopic(row)"
                >
                  批准
                </el-button>
                <el-button
                  v-if="row.status === 'pending'"
                  size="small"
                  @click="rejectTopic(row)"
                >
                  拒绝
                </el-button>
                <el-button
                  v-if="row.status === 'approved'"
                  type="success"
                  size="small"
                  @click="createArticle(row)"
                >
                  创作
                </el-button>
                <el-button
                  type="danger"
                  size="small"
                  @click="deleteTopic(row)"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>
    </el-tabs>

    <!-- 工作流选择对话框 -->
    <el-dialog
      v-model="workflowDialogVisible"
      title="选择工作流"
      width="500px"
    >
      <el-form>
        <el-form-item label="工作流">
          <el-select v-model="selectedWorkflow" placeholder="选择工作流" style="width: 100%">
            <el-option
              v-for="wf in workflows"
              :key="wf.id"
              :label="wf.name"
              :value="wf.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="workflowDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmCreateArticle" :loading="creating">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { Refresh, MagicStick, Search } from '@element-plus/icons-vue';
import * as topicApi from '../../api/topic';
import * as workflowApi from '../../api/workflow';
import { Topic, HotTrend } from '../../types';

const router = useRouter();

// 状态
const activeTab = ref('ai');
const loadingHot = ref(false);
const generating = ref(false);
const loadingLibrary = ref(false);
const creating = ref(false);

// 数据
const aiTopics = ref<Topic[]>([]);
const hotTrends = ref<HotTrend[]>([]);
const libraryTopics = ref<Topic[]>([]);
const workflows = ref<any[]>([]);
const selectedPlatforms = ref(['weibo', 'zhihu', 'baidu']);
const libraryFilter = ref('');
const searchKeyword = ref('');

// 对话框
const workflowDialogVisible = ref(false);
const selectedWorkflow = ref('');
const selectedTopic = ref<Topic | null>(null);

// 名称映射
const platformNames: Record<string, string> = {
  weibo: '微博',
  zhihu: '知乎',
  baidu: '百度',
  toutiao: '今日头条',
};

const sourceNames: Record<string, string> = {
  ai: 'AI生成',
  hot: '热点',
  manual: '手动',
  rss: 'RSS',
};

const statusNames: Record<string, string> = {
  pending: '待处理',
  approved: '已批准',
  rejected: '已拒绝',
  processing: '创作中',
  published: '已发布',
  archived: '已归档',
};

// 计算属性
const groupedHotTrends = computed(() => {
  const groups: Record<string, HotTrend[]> = {};
  for (const trend of hotTrends.value) {
    if (!groups[trend.platform]) {
      groups[trend.platform] = [];
    }
    groups[trend.platform].push(trend);
  }
  return groups;
});

// 方法
function formatHotValue(value: number): string {
  if (value >= 10000) {
    return (value / 10000).toFixed(1) + '万';
  }
  return value.toString();
}

function formatDate(date: string): string {
  return new Date(date).toLocaleString('zh-CN');
}

function getScoreType(score: number): string {
  if (score >= 80) return 'danger';
  if (score >= 70) return 'warning';
  if (score >= 60) return 'success';
  return 'info';
}

function getStatusType(status: string): string {
  const map: Record<string, string> = {
    pending: 'info',
    approved: 'success',
    rejected: 'danger',
    processing: 'warning',
    published: 'success',
    archived: 'info',
  };
  return map[status] || 'info';
}

async function refreshHotTrends() {
  loadingHot.value = true;
  try {
    hotTrends.value = await topicApi.FetchHotTrendsRealtime(
      selectedPlatforms.value,
      20
    );
    ElMessage.success('热点已刷新');
  } catch (error: any) {
    ElMessage.error('刷新失败: ' + error.message);
  } finally {
    loadingHot.value = false;
  }
}

async function fetchHotTrends() {
  await refreshHotTrends();
}

async function generateTopics() {
  generating.value = true;
  try {
    const topics = await topicApi.AIGenerateTopicsFromHot(
      selectedPlatforms.value,
      5
    );
    aiTopics.value = topics;
    activeTab.value = 'ai';
    ElMessage.success(`成功生成 ${topics.length} 个选题`);
  } catch (error: any) {
    ElMessage.error('生成失败: ' + error.message);
  } finally {
    generating.value = false;
  }
}

async function loadLibrary() {
  loadingLibrary.value = true;
  try {
    libraryTopics.value = await topicApi.ListTopics(libraryFilter.value, 50);
  } catch (error: any) {
    ElMessage.error('加载失败: ' + error.message);
  } finally {
    loadingLibrary.value = false;
  }
}

async function searchTopics() {
  if (!searchKeyword.value) {
    await loadLibrary();
    return;
  }
  loadingLibrary.value = true;
  try {
    libraryTopics.value = await topicApi.SearchTopics(searchKeyword.value);
  } catch (error: any) {
    ElMessage.error('搜索失败: ' + error.message);
  } finally {
    loadingLibrary.value = false;
  }
}

async function approveTopic(topic: Topic) {
  try {
    await topicApi.ApproveTopic(topic.id);
    ElMessage.success('已批准');
    await loadLibrary();
  } catch (error: any) {
    ElMessage.error('操作失败: ' + error.message);
  }
}

async function rejectTopic(topic: Topic) {
  try {
    await topicApi.RejectTopic(topic.id);
    ElMessage.success('已拒绝');
    await loadLibrary();
  } catch (error: any) {
    ElMessage.error('操作失败: ' + error.message);
  }
}

async function deleteTopic(topic: Topic) {
  try {
    await ElMessageBox.confirm('确定删除该选题吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    });
    await topicApi.DeleteTopic(topic.id);
    ElMessage.success('已删除');
    await loadLibrary();
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败: ' + error.message);
    }
  }
}

async function saveToLibrary(topic: Topic) {
  try {
    await topicApi.CreateTopic(topic);
    ElMessage.success('已保存到选题库');
    await loadLibrary();
  } catch (error: any) {
    ElMessage.error('保存失败: ' + error.message);
  }
}

async function createArticle(topic: Topic) {
  selectedTopic.value = topic;
  
  // 加载工作流列表
  try {
    workflows.value = await workflowApi.ListWorkflows();
    if (workflows.value.length === 0) {
      ElMessage.warning('请先创建工作流');
      return;
    }
    selectedWorkflow.value = workflows.value[0].id;
    workflowDialogVisible.value = true;
  } catch (error: any) {
    ElMessage.error('加载工作流失败: ' + error.message);
  }
}

async function confirmCreateArticle() {
  if (!selectedTopic.value || !selectedWorkflow.value) return;
  
  creating.value = true;
  try {
    const result = await topicApi.CreateArticleFromTopic(
      selectedTopic.value.id,
      selectedWorkflow.value
    );
    ElMessage.success('已创建文章并启动工作流');
    workflowDialogVisible.value = false;
    
    // 跳转到编辑器
    router.push(`/editor/${result.article.id}`);
  } catch (error: any) {
    ElMessage.error('创建失败: ' + error.message);
  } finally {
    creating.value = false;
  }
}

function openTrendUrl(url: string) {
  if (url) {
    window.open(url, '_blank');
  }
}

// 初始化
onMounted(async () => {
  await refreshHotTrends();
  await loadLibrary();
  
  // 加载平台列表
  try {
    const platforms = await topicApi.GetHotPlatforms();
    if (platforms.length > 0) {
      selectedPlatforms.value = platforms.slice(0, 3);
    }
  } catch (error) {
    // 使用默认值
  }
});
</script>

<style scoped lang="scss">
.topic-center {
  padding: 20px;
  height: 100%;
  overflow-y: auto;
  background: #f5f7fa;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);

  .title-section {
    h2 {
      margin: 0 0 8px 0;
      font-size: 24px;
      color: #303133;
    }
    .subtitle {
      margin: 0;
      color: #909399;
      font-size: 14px;
    }
  }

  .actions {
    display: flex;
    gap: 12px;
  }
}

.topic-tabs {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
}

// AI推荐
.topic-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.topic-card {
  transition: all 0.3s;
  
  &.high-score {
    border: 2px solid #f56c6c;
  }
  
  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  }

  .topic-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
    
    .score {
      font-size: 18px;
      font-weight: bold;
      color: #f56c6c;
    }
  }

  .topic-title {
    margin: 0 0 12px 0;
    font-size: 16px;
    line-height: 1.5;
    color: #303133;
  }

  .topic-reason {
    margin: 0 0 12px 0;
    font-size: 13px;
    color: #606266;
    line-height: 1.6;
  }

  .topic-angles {
    margin-bottom: 12px;
    
    h5 {
      margin: 0 0 8px 0;
      font-size: 12px;
      color: #909399;
    }
    
    .angle-tags {
      display: flex;
      flex-wrap: wrap;
      gap: 6px;
    }
  }

  .topic-keywords {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    margin-bottom: 16px;
  }

  .topic-actions {
    display: flex;
    gap: 8px;
    justify-content: flex-end;
  }
}

// 热点趋势
.hot-trends {
  .platform-filters {
    margin-bottom: 20px;
  }

  .platform-groups {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 20px;
  }

  .platform-section {
    background: #f5f7fa;
    border-radius: 8px;
    padding: 16px;

    .platform-title {
      margin: 0 0 12px 0;
      font-size: 16px;
      color: #303133;
      display: flex;
      align-items: center;
      gap: 8px;

      .count {
        font-size: 12px;
        color: #909399;
        font-weight: normal;
      }
    }
  }

  .trend-items {
    .trend-item {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 10px 12px;
      background: white;
      border-radius: 6px;
      margin-bottom: 8px;
      cursor: pointer;
      transition: all 0.2s;

      &:hover {
        background: #ecf5ff;
      }

      .rank {
        width: 24px;
        height: 24px;
        display: flex;
        align-items: center;
        justify-content: center;
        background: #e4e7ed;
        border-radius: 4px;
        font-size: 12px;
        font-weight: bold;
        color: #606266;
        flex-shrink: 0;

        &.top3 {
          background: #f56c6c;
          color: white;
        }
      }

      .title {
        flex: 1;
        font-size: 14px;
        color: #303133;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .hot-value {
        font-size: 12px;
        color: #f56c6c;
      }
    }
  }
}

// 选题库
.topic-library {
  .filter-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }

  .topic-title-cell {
    .title {
      display: block;
      margin-bottom: 4px;
    }
    .tags {
      display: flex;
      gap: 4px;
      flex-wrap: wrap;
    }
  }
}

.empty-state {
  padding: 60px 0;
}
</style>
