# quack
Tool for create goose migration files according gorm models and database state

----------------
### Restrictions

#### Embed field
 - only anonymous declaration of embed struct

#### Relation fields
 - only explicit declaration for reference fields with tag contains foreignKey attribute

