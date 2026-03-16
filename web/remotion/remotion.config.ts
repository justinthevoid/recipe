import { Config } from '@remotion/cli/config';
import { enableTailwind } from '@remotion/tailwind-v4';

Config.setStillImageFormat('png');
Config.setOverwriteOutput(true);
Config.overrideWebpackConfig(enableTailwind);
