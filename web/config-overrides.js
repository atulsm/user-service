module.exports = function override(config, env) {
  // Disable Fast Refresh
  if (env === 'development') {
    config.module.rules = config.module.rules.map(rule => {
      if (rule.oneOf) {
        rule.oneOf = rule.oneOf.map(oneOfRule => {
          if (oneOfRule.loader && oneOfRule.loader.includes('babel-loader')) {
            oneOfRule.options = {
              ...oneOfRule.options,
              plugins: (oneOfRule.options.plugins || []).filter(
                plugin => !plugin.includes('react-refresh')
              ),
            };
          }
          return oneOfRule;
        });
      }
      return rule;
    });
  }
  return config;
}; 