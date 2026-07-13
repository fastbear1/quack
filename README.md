# quack
Tool for create goose migration files according gorm models and database schema

----------------
### Restrictions

#### Embed field
 - only anonymous declaration of embed struct

#### Relation fields
 - only explicit declaration for reference fields with tag contains foreignKey attribute

#### Alter column
 - only data, nullable and default value checking for altering column

#### parsing Index tags
 - only unique, index type, expression and column list are used for creaeting indices. Same params used to check index state(alter index)
