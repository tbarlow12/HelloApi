ALTER TABLE [dbo].[Product_Category_Rel]
	ADD CONSTRAINT [FK_Product_Category_Product]
	FOREIGN KEY (ProductId)
	REFERENCES [dbo].[Product] (Id)
