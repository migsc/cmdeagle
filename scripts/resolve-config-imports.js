const fs = require("fs");
const path = require("path");
const YAML = require("yaml");

async function main() {
  const cwd = process.cwd();
  const configFilePath = path.join(cwd, ".cmd.yaml");
  const configObject = await loadYAMLFile(configFilePath);
  await resolveActionPaths(configObject);
  console.log(YAML.stringify(configObject, null, 2));
}

main();

async function loadYAMLFile(filePath) {
  const content = fs.readFileSync(filePath, "utf-8");
  const parsed = YAML.parseDocument(content);
  const parsedObj = parsed.toJS();
  await processImports(parsedObj);
  return parsedObj;
}

async function processImports(obj) {
  if (!obj || typeof obj !== 'object') return;

  if (Array.isArray(obj)) {
    for (let i = 0; i < obj.length; i++) {
      const item = obj[i];
      if (Array.isArray(item)) {
        obj.splice(i, 1, ...item.flat());
        i--;
        continue;
      }
      
      if (item && typeof item === 'object') {
        if (item.imports) {
          const imports = Array.isArray(item.imports) ? item.imports : [item.imports];
          const importedItems = [];
          for (const importPath of imports) {
            let importedData;
            if (isURL(importPath)) {
              importedData = await resolveUrlImport(importPath);
            } else {
              importedData = await resolveImportPath(importPath);
            }
            importedItems.push(importedData);
          }
          obj.splice(i, 1, ...importedItems.flat());
          i += importedItems.length - 1;
        } else {
          await processImports(item);
        }
      }
    }
    
    const flattened = obj.flat();
    obj.length = 0;
    obj.push(...flattened);
    return;
  }

  for (const [_, value] of Object.entries(obj)) {
    if (value && typeof value === 'object') {
      await processImports(value);
    }
  }

  if (obj.imports) {
    const imports = Array.isArray(obj.imports) ? obj.imports : [obj.imports];
    for (const importPath of imports) {
      let importedData;
      if (isURL(importPath)) {
        importedData = await resolveUrlImport(importPath);
      } else {
        importedData = await resolveImportPath(importPath);
      }
      Object.assign(obj, importedData);
    }
    delete obj.imports;
  }
}

function isURL(path) {
  try {
    new URL(path);
    return true;
  } catch {
    return false;
  }
}

async function resolveImportPath(filePath) {
  const content = fs.readFileSync(filePath, "utf-8");
  const parsed = YAML.parse(content);
  const normalized = YAML.parseDocument(YAML.stringify(parsed));
  const parsedObj = normalized.toJS();
  await processImports(parsedObj);
  return parsedObj;
}

async function resolveUrlImport(urlValue) {
  try {
    const response = await fetch(urlValue);
    if (!response.ok) {
      throw new Error(`Failed to fetch URL: ${urlValue}`);
    }
    const content = await response.text();
    const parsed = YAML.parse(content);
    const normalized = YAML.parseDocument(YAML.stringify(parsed));
    const parsedObj = normalized.toJS();
    await processImports(parsedObj);
    return parsedObj;
  } catch (error) {
    console.error(`Error importing from URL ${urlValue}:`, error);
    throw error;
  }
}

async function resolveActionPaths(config) {
  let basePath = config.from;

  for (const command of config.commands ?? []) {
    basePath = command.from ?? config.from;

    if (typeof command.action === "string") {
      command.action = path.resolve(command.from ?? config.from, command.action);
    }

    if (command.commands) {
      await resolveActionPaths(command);
    }
  }
  return config;
} 