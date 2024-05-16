## checker 說明
### 使用方法

把要檢查的網址列寫在 a 檔案，一行一個網域
建議是完整網域，有www就加，子域名也可以

#### from.txt

 www.abc.com

 www.def.com

執行完畢後如果有有效結果，就會叫你輸入檔案存檔
預設檔名會是 result.txt，一樣是一個結果一行

### 修改程式碼
檢查的路徑跟檢查的內容要自己改，check_path 要包含網域後的完整路徑、check_string 是完整比對，不能少

``` go
const check_path = "__PATH__"
const check_string = "__CHECK_STR__"
```