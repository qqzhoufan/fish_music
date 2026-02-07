-- 更新所有歌曲的国家代码，基于语言字段
UPDATE songs SET country_code = 'CN' WHERE language = '华语';
UPDATE songs SET country_code = 'US' WHERE language = '英语';
UPDATE songs SET country_code = 'JP' WHERE language = '日语';
UPDATE songs SET country_code = 'KR' WHERE language = '韩语';
UPDATE songs SET country_code = 'FR' WHERE language = '法语';
UPDATE songs SET country_code = 'DE' WHERE language = '德语';
UPDATE songs SET country_code = 'ES' WHERE language = '西班牙语';
UPDATE songs SET country_code = 'RU' WHERE language = '俄语';
UPDATE songs SET country_code = 'IT' WHERE language = '意大利语';
UPDATE songs SET country_code = 'BR' WHERE language = '葡萄牙语';
UPDATE songs SET country_code = 'TH' WHERE language = '泰语';
UPDATE songs SET country_code = 'VN' WHERE language = '越南语';
UPDATE songs SET country_code = 'ID' WHERE language = '印尼语';
UPDATE songs SET country_code = 'MY' WHERE language = '马来语';
UPDATE songs SET country_code = 'IN' WHERE language = 'Hindi';
UPDATE songs SET country_code = 'PH' WHERE language = 'Tagalog';
UPDATE songs SET country_code = 'US' WHERE language = '其他';
