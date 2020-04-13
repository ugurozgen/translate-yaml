# tly 
it translates values of yaml to another language by using google translate api

for example; english to turkish yaml tanslation

filename: ```caption.yml```
```
header:
    title: hello world!
```

turns...

filename: ```caption.yml-tr```
```    
header:
    title: merhaba dunya!  
```

### Build
./build.sh 

### Run 
```
$ tly -from en -to tr -f ./test.yml

$ tly -from en -to tr -f .
```

### Arguments
```
-from //source language
-to   //target language
-f    //file or folder path
``` 
