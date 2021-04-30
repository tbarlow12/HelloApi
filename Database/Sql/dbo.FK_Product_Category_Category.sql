ALTER TABLE [dbo].[Product_Category_Rel]
	ADD CONSTRAINT [FK_Product_Category_Category]
	FOREIGN KEY (CategoryId)
	REFERENCES [dbo].[Category] (Id)
