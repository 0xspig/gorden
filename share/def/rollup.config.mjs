// rollup.config.mjs
import commonjs from "@rollup/plugin-commonjs";
import nodeResolve from "@rollup/plugin-node-resolve";
import scss from "rollup-plugin-scss"

export default {
	input: 'src/main.js',
	output: {
		file: 'static/main.js',
		format: 'es'
	},
	plugins:[
		commonjs(),
		nodeResolve(),
		scss({ fileName: 'main.css' })
	]
};