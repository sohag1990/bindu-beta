# bindu-beta
Simple api development go framework top on gin

bindu new  --app Blank --port 8080 #true

bindu new Hello –app Blank –db 'AdapterName:HostName:Port:DbName:DbUserName:DbPass' –port 8080

bindu new Hello –app blank –port 9999 –db Mysql:Localhost:3306:blog:root



bindu generate Help #true

bindu generate Scaffold User Email:String UserName:String Phone:String --hasOne Profile 

bindu generate Model Profile FirstName:String LastName:String Image:String  

bindu generate Controller Profile FirstName:String LastName:String Image:String  

bindu generate Scaffold Item Location:String Type:String Image:String Title:String Date:String 

bindu generate Model User --hasMany Item #false

bindu create User Username:Sohag Password:11111 Role:Admin 
bindu generate Routes User --middleware auth --group v1 
bindu generate Routes Item --middleware auth --group v1 
bindu generate Controller User

bindu create Policy Alice:P Sub:Admin Obj:/Api/* Act:* 

bindu generate Scaffold Blog Title:String Description:String --hasMany Comment --belongsTo User 

bindu generate Scaffold Comment Subject:String Body:String 


bindu generate Scaffold Contact FirstName:String LastName:String Subject:String Description:String

bindu generate Model Item Description:String 

bindu generate Model Item Mobile:String 

bindu db Migrate


bindu create Policy Alice:P Sub:Subscriber Obj:/Api/V1/* Act:* 

bindu create User Username:Sohag Password:11111 Role:Admin 


// Swagger
bindu add swagger
bindu add swagger –skip // to skip download lib and faster execution

// DB
bindu db migrate –db Dbname –table tableName

//fix
bindu fix
bindu fix import –oldPath github.com/old/path

// RelationShip:

hasOne /
hasMany // model field type fixed korte hobe
manyToMany // self Referencing manytomany add korte hobe
belongsTo // association key fix korte hobe
hasOneThrough
hasManyThrough


Draft

bindu generate model User name:string, title:string_512 bio:text –hasMany posts:postID –hasOne address-profile

bindu generate controller User –methods all

bindu generate controller Profile –methods get-post-put-delete

bindu generate scaffold User name:string, title:string –hasMany posts –hasOne address

bindu generate routes User

// update

bindu update model User name:string, title:string –hasMany posts –hasOne address

bindu update controller User

bindu update scaffold User name:string, title:string –hasMany posts –hasOne address

bindu update routes User

// important
bindu generate routes User -g=v1 -m=auth



// Create

bindu create user username:Sohag, password:11111 role:admin // admin must be in smallcase

bindu create Policy alice:p, sub:Admin, obj:/api/*, act:*

bindu create Policy alice:p, sub:editor, obj:/api/v2/ping/*, act:*

// Modify
bindu modify model User name –addGorm uniqe, not null, type:text, type:varchar(50) –removeGorm   index

// History 
bindu history // to see all command executed
// install
bindu install

// serve
bindu serve –port 8080 –host localhost




