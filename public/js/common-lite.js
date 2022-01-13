// leanote 通用方法

/**
 * 统计字数（中文按单字算，英文按单词算)
 */
 function calcWords(str) {
    sLen = 0;
    try{
        str = str.replace(/<[^>]+>/g,"");        //过滤所有的HTML标签
        //先将回车换行符做特殊处理
        str = str.replace(/(\r\n+|\s+|　+)/g,"龘");
        //处理英文字符数字，连续字母、数字、英文符号视为一个单词
        str = str.replace(/[\x00-\xff]/g,"m");  
        //合并字符m，连续字母、数字、英文符号视为一个单词
        str = str.replace(/m+/g,"*");
        //去掉回车换行符
        str = str.replace(/龘+/g,"");
        //返回字数
        sLen = str.length;
    }catch(e){
    }
    return sLen;
}