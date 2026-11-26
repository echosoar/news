const fs = require("node:fs");
const path = require("node:path");

const ROOT = process.cwd();
const TAGS_DIR = path.join(ROOT, "tags");
const GROUPS_DIR = path.join(ROOT, "groups");
const DATE_FILE_RE = /^\d{4}-\d{2}-\d{2}\.json$/;
const PAGE_SIZE = 200;
const SIZE_LIMIT = 5;
const STATS_DAY_COUNT = 100;

fs.mkdirSync(TAGS_DIR, { recursive: true });
fs.mkdirSync(GROUPS_DIR, { recursive: true });

const generatedAt = new Date().toISOString().replace(/\.\d{3}Z$/, "Z");

function readJson(filePath) {
  return JSON.parse(fs.readFileSync(filePath, "utf8"));
}

function uniqueStrings(values) {
  const result = [];
  for (const value of values) {
    const normalized = String(value).trim();
    if (normalized && !result.includes(normalized)) {
      result.push(normalized);
    }
  }
  return result;
}

function escapeRegExp(value) {
  return value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}

const groupDefinitions = [];
const aliasDefinitions = [];
const groupConfigPath = path.join(ROOT, "group.json");

if (fs.existsSync(groupConfigPath)) {
  try {
    const groupPayload = readJson(groupConfigPath);
    if (groupPayload.alias && typeof groupPayload.alias === "object" && !Array.isArray(groupPayload.alias)) {
      for (const [rawPattern, rawTargets] of Object.entries(groupPayload.alias)) {
        const pattern = String(rawPattern).trim();
        if (!pattern) continue;

        const targets = Array.isArray(rawTargets) ? rawTargets : [rawTargets];
        const normalizedTargets = uniqueStrings(targets);
        const source = `^${escapeRegExp(pattern).replace(/\\\*/g, ".*")}$`;
        aliasDefinitions.push({ pattern: new RegExp(source), targets: normalizedTargets });
      }
    }

    if (Array.isArray(groupPayload.groups)) {
      for (const rawGroup of groupPayload.groups) {
        if (!rawGroup || typeof rawGroup !== "object" || Array.isArray(rawGroup)) continue;

        const group = String(rawGroup.group ?? "").trim();
        if (!group) continue;

        const tags = Array.isArray(rawGroup.tags) ? uniqueStrings(rawGroup.tags) : [];
        groupDefinitions.push({ group, tags });
      }
    }
  } catch {
    // An invalid config leaves both indexes empty, matching the previous generator.
  }
}

function resolveAliases(tag) {
  const aliases = [];
  let matched = false;

  for (const definition of aliasDefinitions) {
    if (!definition.pattern.test(tag)) continue;
    matched = true;
    for (const target of definition.targets) {
      if (!aliases.includes(target)) aliases.push(target);
    }
  }

  return matched ? aliases : [tag];
}

for (const definition of groupDefinitions) {
  const resolvedTags = [];
  for (const tag of definition.tags) {
    for (const resolvedTag of resolveAliases(tag)) {
      if (!resolvedTags.includes(resolvedTag)) resolvedTags.push(resolvedTag);
    }
  }
  definition.tags = resolvedTags;
}

function stableSortValue(value) {
  if (Array.isArray(value)) return value.map(stableSortValue);
  if (value && typeof value === "object") {
    return Object.fromEntries(
      Object.keys(value)
        .sort()
        .map((key) => [key, stableSortValue(value[key])]),
    );
  }
  return value;
}

function buildDailyStats(entries) {
  const endDate = new Date(`${generatedAt.slice(0, 10)}T00:00:00Z`);
  const startDate = new Date(endDate);
  startDate.setUTCDate(startDate.getUTCDate() - STATS_DAY_COUNT + 1);
  const counts = new Map();

  for (const entry of entries) {
    const dateString = String(entry._sort_date ?? "");
    const date = new Date(`${dateString}T00:00:00Z`);
    if (Number.isNaN(date.getTime()) || date < startDate || date > endDate) continue;
    counts.set(dateString, (counts.get(dateString) || 0) + 1);
  }

  const stats = [];
  const currentDate = new Date(startDate);
  for (let offset = 0; offset < STATS_DAY_COUNT; offset += 1) {
    const date = currentDate.toISOString().slice(0, 10);
    stats.push({ date, count: counts.get(date) || 0 });
    currentDate.setUTCDate(currentDate.getUTCDate() + 1);
  }
  return stats;
}

const tagEntries = new Map();
const groupEntries = new Map(groupDefinitions.map(({ group }) => [group, []]));
const groupSeenKeys = new Map(groupDefinitions.map(({ group }) => [group, new Set()]));
const generatedFiles = new Set();

const dateFiles = fs
  .readdirSync(ROOT, { withFileTypes: true })
  .filter((entry) => entry.isFile() && DATE_FILE_RE.test(entry.name))
  .map((entry) => entry.name)
  .sort()
  .reverse();

for (const fileName of dateFiles) {
  let payload;
  try {
    payload = readJson(path.join(ROOT, fileName));
  } catch {
    continue;
  }

  const dateValue = payload && typeof payload === "object" ? payload.date : undefined;
  const summaries = payload && typeof payload === "object" ? payload.summary : undefined;
  if (!Array.isArray(summaries)) continue;

  for (const item of summaries) {
    if (!item || typeof item !== "object" || Array.isArray(item)) continue;

    let tagList = [];
    if (Array.isArray(item.tags)) {
      for (const rawTag of item.tags) {
        const normalizedTag = String(rawTag).trim();
        if (!normalizedTag) continue;
        for (const resolvedTag of resolveAliases(normalizedTag)) {
          if (!tagList.includes(resolvedTag)) tagList.push(resolvedTag);
        }
      }
    }
    if (!Array.isArray(item.tags) || item.tags.length === 0) tagList = ["未分类"];

    let sortDate = path.basename(fileName, ".json");
    if (typeof dateValue === "string" && DATE_FILE_RE.test(`${dateValue}.json`)) {
      sortDate = dateValue;
    }

    const normalized = structuredClone(item);
    normalized.tags = tagList;
    normalized.date = dateValue;
    normalized._sort_date = sortDate;

    const itemKey = JSON.stringify(
      stableSortValue({
        date: dateValue,
        title: normalized.title,
        summary: normalized.summary,
        links: item.links,
        tags: tagList,
      }),
    );

    delete normalized.links;

    for (const tag of tagList) {
      if (!tagEntries.has(tag)) tagEntries.set(tag, []);
      tagEntries.get(tag).push(structuredClone(normalized));
    }

    const itemTags = new Set(tagList);
    for (const definition of groupDefinitions) {
      if (!definition.tags.some((tag) => itemTags.has(tag))) continue;

      const seenKeys = groupSeenKeys.get(definition.group);
      if (seenKeys.has(itemKey)) continue;

      seenKeys.add(itemKey);
      groupEntries.get(definition.group).push(structuredClone(normalized));
    }
  }
}

function writePaginatedFiles(directory, name, kind, entries) {
  entries.sort((left, right) => String(right._sort_date ?? "").localeCompare(String(left._sort_date ?? "")));
  const stats = buildDailyStats(entries);
  for (const entry of entries) delete entry._sort_date;

  const pageCount = Math.max(1, Math.ceil(entries.length / PAGE_SIZE));
  for (let page = 1; page <= pageCount; page += 1) {
    const suffix = page === 1 ? "" : `-${page}`;
    const relativeName = `${directory}/${name}${suffix}.json`;
    generatedFiles.add(relativeName);

    const output = {
      [kind]: name,
      generated_at: generatedAt,
      page,
      page_count: pageCount,
      count: entries.length,
    };
    if (page === 1) output.stats = stats;
    output.summary = entries.slice((page - 1) * PAGE_SIZE, page * PAGE_SIZE);

    fs.writeFileSync(
      path.join(ROOT, relativeName),
      `${JSON.stringify(output, null, 2)}\n`,
      "utf8",
    );
  }

  return `${directory}/${name}.json`;
}

const tagItems = [];
for (const [tag, entries] of tagEntries) {
  if (entries.length < SIZE_LIMIT) continue;

  const file = writePaginatedFiles("tags", tag, "tag", entries);
  tagItems.push({ tag, count: entries.length, file });
}

for (const [group, entries] of groupEntries) {
  writePaginatedFiles("groups", group, "group", entries);
}

for (const indexDir of [TAGS_DIR, GROUPS_DIR]) {
  for (const entry of fs.readdirSync(indexDir, { withFileTypes: true })) {
    if (!entry.isFile() || !entry.name.endsWith(".json")) continue;
    const relativeName = `${path.basename(indexDir)}/${entry.name}`;
    if (!generatedFiles.has(relativeName)) fs.unlinkSync(path.join(indexDir, entry.name));
  }
}

tagItems.sort((left, right) => right.count - left.count || left.tag.localeCompare(right.tag));
fs.writeFileSync(
  path.join(ROOT, "tags.json"),
  `${JSON.stringify({ generated_at: generatedAt, tags: tagItems }, null, 2)}\n`,
  "utf8",
);
