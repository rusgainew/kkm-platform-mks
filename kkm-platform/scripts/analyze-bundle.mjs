#!/usr/bin/env node

/**
 * Bundle Size Analyzer для Next.js проекта
 * Анализирует размер JS bundle и выдает рекомендации
 * 
 * Usage: npx node scripts/analyze-bundle.mjs
 */

import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const projectRoot = path.join(__dirname, '..');
const buildDir = path.join(projectRoot, '.next');

interface FileInfo {
  name: string;
  size: number;
  gzipSize?: number;
}

function formatBytes(bytes) {
  const kb = bytes / 1024;
  const mb = kb / 1024;
  
  if (mb > 1) return `${mb.toFixed(2)} MB`;
  if (kb > 1) return `${kb.toFixed(2)} KB`;
  return `${bytes} B`;
}

function analyzeBundle() {
  console.log('📊 Bundle Size Analysis\n');
  console.log('='.repeat(60));

  const analysisDir = path.join(buildDir, 'static', 'chunks');
  
  if (!fs.existsSync(analysisDir)) {
    console.log('❌ Build directory not found. Run "npm run build" first.\n');
    return;
  }

  const files: FileInfo[] = [];
  let totalSize = 0;
  let maxSize = 0;
  let maxFile = '';

  // Сканируем файлы
  const items = fs.readdirSync(analysisDir);
  items.forEach(file => {
    if (file.endsWith('.js')) {
      const filePath = path.join(analysisDir, file);
      const stats = fs.statSync(filePath);
      const size = stats.size;
      
      files.push({
        name: file,
        size: size,
      });

      totalSize += size;

      if (size > maxSize) {
        maxSize = size;
        maxFile = file;
      }
    }
  });

  // Сортируем по размеру
  files.sort((a, b) => b.size - a.size);

  // Выводим топ 10 самых больших файлов
  console.log('\n📈 Top 10 Largest Chunks:\n');
  files.slice(0, 10).forEach((file, i) => {
    const percentage = ((file.size / totalSize) * 100).toFixed(1);
    const bar = '█'.repeat(Math.round((file.size / maxSize) * 40));
    console.log(`${(i + 1).toString().padEnd(2)} ${file.name.padEnd(40)} ${formatBytes(file.size).padEnd(12)} ${bar} ${percentage}%`);
  });

  console.log('\n' + '='.repeat(60));
  console.log(`\n📦 Total Bundle Size: ${formatBytes(totalSize)}`);
  console.log(`📁 Number of chunks: ${files.length}`);
  console.log(`🔴 Largest chunk: ${maxFile} (${formatBytes(maxSize)})`);

  // Анализ и рекомендации
  console.log('\n💡 Recommendations:\n');

  const avgSize = totalSize / files.length;
  const largeFiles = files.filter(f => f.size > avgSize * 2);

  if (largeFiles.length > 0) {
    console.log('⚠️  Found large chunks that could be code-split further:');
    largeFiles.forEach(file => {
      console.log(`   - ${file.name} (${formatBytes(file.size)})`);
    });
    console.log('   Consider using dynamic imports or lazy loading.\n');
  }

  if (totalSize > 500 * 1024) {
    console.log('⚠️  Total bundle size is larger than recommended (>500KB).');
    console.log('   Consider implementing:');
    console.log('   - Code splitting for heavy components');
    console.log('   - Dynamic imports for routes');
    console.log('   - Tree shaking and minification\n');
  } else {
    console.log(`✅ Bundle size is healthy (${formatBytes(totalSize)}).\n`);
  }

  // Проверка специфических библиотек
  console.log('📚 Detected libraries in bundle:');
  const libPatterns = [
    { name: 'React', pattern: /react/ },
    { name: 'Next.js', pattern: /next/ },
    { name: 'recharts', pattern: /recharts/ },
    { name: 'Zustand', pattern: /zustand/ },
    { name: 'TanStack Query', pattern: /@tanstack|react-query/ },
  ];

  libPatterns.forEach(lib => {
    const found = files.some(f => lib.pattern.test(f.name));
    console.log(`   ${found ? '✅' : '⭕'} ${lib.name}`);
  });

  console.log('\n' + '='.repeat(60) + '\n');
}

analyzeBundle();
