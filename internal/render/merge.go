package render

func DeepMerge(base, over map[string]any) map[string]any {
	out := make(map[string]any, len(base))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range over {
		if existing, ok := out[k]; ok {
			if baseMap, baseOk := existing.(map[string]any); baseOk {
				if overMap, overOk := v.(map[string]any); overOk {
					out[k] = DeepMerge(baseMap, overMap)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}
