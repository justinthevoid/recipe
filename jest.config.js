// Jest configuration for Recipe web interface tests
// Story: 11-1-css-filter-mapping
export default {
    testEnvironment: 'jsdom',
    transform: {
        '^.+\\.js$': 'babel-jest',
    },
    moduleFileExtensions: ['js'],
    testMatch: ['**/__tests__/**/*.test.js'],
    collectCoverageFrom: [
        'web/static/**/*.js',
        '!web/static/**/*.test.js',
        '!web/static/bundle.min.js',
        '!web/static/wasm_exec.js',
    ],
    coverageThreshold: {
        'web/static/preview.js': {
            functions: 85,
            lines: 90,
            statements: 90,
        },
    },
    coverageReporters: ['text', 'lcov', 'html'],
};
