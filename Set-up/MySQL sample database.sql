CREATE DATABASE  IF NOT EXISTS `ProjectGoLive` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `ProjectGoLive`;
-- MySQL dump 10.13  Distrib 8.0.22, for macos10.15 (x86_64)
--
-- Host: localhost    Database: ProjectGoLive
-- ------------------------------------------------------
-- Server version	8.0.22

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `Orders`
--

DROP TABLE IF EXISTS `Orders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Orders` (
  `OrderID` int NOT NULL AUTO_INCREMENT,
  `UserID` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `ProductID` int NOT NULL,
  `Qty` int NOT NULL,
  `SellerID` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `Status` int NOT NULL,
  `Created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Modified` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`OrderID`),
  KEY `UserID` (`UserID`),
  KEY `ProductID` (`ProductID`),
  KEY `SellerID` (`SellerID`),
  CONSTRAINT `orders_ibfk_1` FOREIGN KEY (`UserID`) REFERENCES `User` (`UserID`),
  CONSTRAINT `orders_ibfk_2` FOREIGN KEY (`ProductID`) REFERENCES `Product` (`ProductID`),
  CONSTRAINT `orders_ibfk_3` FOREIGN KEY (`SellerID`) REFERENCES `User` (`UserID`)
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Orders`
--

LOCK TABLES `Orders` WRITE;
/*!40000 ALTER TABLE `Orders` DISABLE KEYS */;
INSERT INTO `Orders` VALUES (19,'ongryan123',36,1,'ahmadmuhammad',0,'2021-05-13 04:37:25','2021-05-13 04:37:25'),(27,'ongryan123',38,6,'leematthew',1,'2021-05-13 12:33:46','2021-05-13 12:33:46'),(28,'ongryan123',52,1,'leematthew',2,'2021-05-13 12:33:46','2021-05-13 12:33:46'),(29,'ongryan123',39,5,'leematthew',1,'2021-05-13 13:12:47','2021-05-13 13:12:47'),(30,'ongryan123',39,5,'leematthew',2,'2021-05-14 02:12:24','2021-05-14 02:12:24'),(31,'ongryan123',52,3,'leematthew',1,'2021-05-14 02:12:24','2021-05-14 02:12:24');
/*!40000 ALTER TABLE `Orders` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Product`
--

DROP TABLE IF EXISTS `Product`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Product` (
  `ProductID` int NOT NULL AUTO_INCREMENT,
  `Name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `Description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `Keyword` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `CategoryID` int NOT NULL,
  `Price` float DEFAULT NULL,
  `DiscountID` int DEFAULT '0',
  `Created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Modified` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `SellerID` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `Inventory` int NOT NULL,
  `Rating` float NOT NULL DEFAULT '0',
  `RatingNum` int NOT NULL DEFAULT '0',
  `UnitSold` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`ProductID`),
  KEY `CategoryID` (`CategoryID`),
  KEY `DiscountID` (`DiscountID`),
  KEY `SellerID` (`SellerID`),
  CONSTRAINT `product_ibfk_1` FOREIGN KEY (`SellerID`) REFERENCES `User` (`UserID`)
) ENGINE=InnoDB AUTO_INCREMENT=68 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Product`
--

LOCK TABLES `Product` WRITE;
/*!40000 ALTER TABLE `Product` DISABLE KEYS */;
INSERT INTO `Product` VALUES (32,'Grass Fed Ribeye Steak (300g)','Our Premium Angus Ribeye Steaks is a restaurant menu staple due to its incredible marbling and tenderness. This cut is suited to both grilling and pan-frying. You can be assured of the food safety for our beef, as it is cut and vacuum packed at the processor in Singapore. This means less contamination risks as third-party handling is significantly reduced, with no-one ‘touching’ your meat until it is opened by you in the comfort of your home. ','Ribeye Steak',2,24.8,1,'2021-05-06 13:24:12','2021-05-06 13:24:12','jasonlim',47,5,15,33),(34,'Fresh Red Snapper ','Grown in high tech farm in Singapore near Changi Point Ferry terminal. Can be cooked in multiple ways. Delivery day 5-7 days.','Snapper Fish',2,35,2,'2021-05-06 13:33:51','2021-05-06 13:33:51','jasonlim',100,4.7,20,50),(35,'Frozen Minced Pork (500g)','500gm packet of Pork Mince. Perfect for all sorts of dishes, burgers, sausages and more.This item is FROZEN and will be delivered FROZEN.  It is OK to keep FROZEN for approx. 6 months, then once defrosted can be kept in the fridge chilled for approx. 2-3 days.','MInced Pork',0,18,0,'2021-05-06 13:33:51','2021-05-06 13:33:51','jasonlim',27,4.9,30,103),(36,'Red Onion (3kg)','Red onion is rich in nutrients, can sterilize, promote digestion, help lower blood pressure and blood lipids, is an indispensable vegetable in many foods.\n\nHigh quality red onions from India with large packaging can meet the needs of large families.','Onion',5,6.3,0,'2021-05-06 13:38:29','2021-05-06 13:38:29','ahmadmuhammad',200,4.5,200,400),(37,'Red Honey Cherry Tomatoes +/-250g','Organic tomatoes grown locally. Sweet flavour, store in cooling places.\r \r Storage 2 weeks.','Tomatoes',4,3,0,'2021-05-06 13:40:22','2021-05-06 13:40:22','jasonlim',20,3,10,100),(38,'The Brewery Bakery: Sticky Bun Beer~ P493','We are a Singaporean brewery that is set out to change the world’s perception on Singaporean craft beer. Dig into this bold barrel-aged imperial stout and experience a pipin’ hot plate of liquid sticky buns in your glass. This barrel-aged imperial stout is layered with pecans, maple syrup, and just a hint of cinnamon — imitating the warming flavors of freshly-baked sticky buns. This bakery is open for business.','Beer',3,26.5,4,'2021-05-06 13:44:52','2021-05-06 13:44:52','leematthew',269,5,20,51),(39,'White Rabbit White Ale Beer ~ P2','White Rabbit White Ale celebrates the tradition of brewing with imagination and creativity and delivers refreshing hints of coriander, juniper berry and bitter orange, blended with a hefty dose unmalted wheat. Light citrus aromas round out a classic, cloudy white ale with just a gentle hint of bitterness.','Beer',3,7,5,'2021-05-06 13:45:47','2021-05-06 13:45:47','leematthew',187,4.8,300,3013),(40,'Rye & Pint. Star Gazin Beer~ P104','Dry-hopped with mainly Galaxy hops, Star Gazin’ presents itself with notes of tropical fruits aroma and flavor, paired with a subtle malt sweetness, making it easy-drinking IPA.','Beer',3,8,0,'2021-05-06 13:46:49','2021-05-06 13:46:49','leematthew',100,4.2,189,389),(41,'Fresh Bananas, 6 lbs. 2/Pack (900-00106)','Deliciously healthy bananas for the office break room! Delivered fresh, this fruit will make any work place sweeter! Perfect on the go snack, and a great way to snack healthy at work!\n\n6 pounds of ripe fresh bananas\nFresh fruit for the office breakroom\nGet this fresh snack for the healthy eaters in your office','Bananas',4,11.49,0,'2021-05-06 13:52:57','2021-05-06 13:52:57','ahmadmuhammad',250,5,10,12),(42,'Fresh Cinnamon Raisin Bagels, 6/Pack (900-00008)','Warm up the toaster for these fresh sliced bagels. Start off your work day the right way with a delicious breakfast of a warm toasted bagel. Breakfast at the office should be quick and taste great, and this 6-Pack of Bagels is perfect for the office break room.\n\nA 6-Pack of fresh Cinnamon bagels\nPre-sliced bagels for a quick breakfast at the office\nDeliciously soft Cinnamon Bagels make a breakfast great','Bagels',1,12,1,'2021-05-06 13:54:31','2021-05-06 13:54:31','ahmadmuhammad',40,4.5,18,32),(43,'Fresh Organic White Rice 1kg','Ingredient List: 100% White Rice. 100% whole-kernel. White rice is the name given to milled rice that has its husk, bran and germ removed. After milling and polishing, the rice is bright, white and shiny. However, important nutrients such as vitamin B1 and B3 are removed. Fresh Rice is sourced from Surin in North east Thailand. Surin is known to produce the best \"Hom Mali\" rice because of its ideal weather conditions and vast fertile land.','Rice',1,15,0,'2021-05-06 13:57:35','2021-05-06 13:57:35','balamuthu',29,4.9,218,344),(44,'Golden Ale - 4 x 330ml','A clean, crisp, richly golden hued ale brewed with premium British malts and specially selected varieties of hops. This ale has a balanced bread and biscuit maltiness accented by smooth bitterness and a mild floral and citrus aroma and flavour.','Beer',3,23.2,0,'2021-05-07 12:22:05','2021-05-07 12:22:05','leematthew',300,4.9,26,100),(45,'[20% OFF] Bohemian Pilsner - 6 x 330ml','A Bohemian Pilsner with light notes of honey and a crisp, wonderful balance of malt and bitterness brought out by our specially selected yeast from a world famous brewery in the Czech Republic.\n\n5.0% abv / 30 IBU','Beer',3,27.84,3,'2021-05-07 12:23:17','2021-05-07 12:23:17','leematthew',154,4.7,30,280),(46,'Seng Choon Lower Cholesterol Eggs - Farm Fresh','UV sterilizedThese fresh eggs are produced by hens fed with a wholesome mix of natural grains. Hence, these eggs are of the highest quality with lower cholesterol.','Eggs',1,2.75,0,'2021-05-07 12:25:40','2021-05-07 12:25:40','balamuthu',2000,4.5,200,211),(47,'Dasoon Premium Fresh Eggs','Eggs are rich and flavorful! Incorporate them into your meals for a healthy and balanced diet.Health benefits:\n• High in protein: Strengthen bone health\n• Contains vitamins & minerals: Enhance brain functions and promote healthy skin\n• Contains choline: Keep heart healthy\nCountry of Origin\nMalaysia\n','Eggs',1,3.75,1,'2021-05-07 12:28:40','2021-05-07 12:28:40','balamuthu',300,4.8,20,140),(48,'Yili Farm Premium Kang Kong','Long hollow stems with pointed, mid-green leaves. Also known as morning glory, water convolvulus, water spinach, swap cabbage, ong choy, hung tsai, rau muong. Yili Premium products are under SFAs GAP (Good Agriculture Practice) & Love SG logo. It has an extremely mild flavour and also slightly sweet. This dark green vegetable is extremely rich in iron and calcium. It also contains Vitamin A, B & C.To prepare, separate stems and leaves. While cooking, add the stems first as leaves wilt quickly.Commonly cooked with Sambal/Belachan chilli/garlic. It has a soft crunchy texture. Best consumed within 3 days upon receiving the item.','Water Spinach',4,2.55,0,'2021-05-07 12:30:57','2021-05-07 12:30:57','ahmadmuhammad',100,4,39,300),(51,'Pasar Seedless Red Watermelon','As its name suggests, this watermelon has much less seeds than other varieties. Its juicy red flesh is also sweeter and easier to consume. Hydrate yourself with a slice of watermelon on a hot day, or blend its flesh for a refreshing drink.','Watermelon',4,10,0,'2021-05-07 12:33:29','2021-05-07 12:33:29','ahmadmuhammad',19,5,188,203),(52,'ROAD HOG SESSION IPA (24Pack)','Delightfully \"hoppy\" creating a crisp finish. Hosting citrus hop notes of juniper berries to a grainy roasted malt, this beer has a rounded bitterness. Available as either a 6 pack, 12 pack or case of 24: bottles are 330mL.','Beer',3,85,1,'2021-05-07 12:36:06','2021-05-07 12:36:06','leematthew',459,4.3,273,496),(60,'STOLEN BOAT SUMMER ALE (12Pack)','A summer sweetness from the malt, with a short bitter finish. This Ale hosts sweet bread notes with hints of honey marmalade and a caramel, fruity hop finish. Available as either a 6 pack, 12 pack or a case of 24: bottles are 330mL.','Beer',3,42,0,'2021-05-07 12:37:28','2021-05-07 12:37:28','leematthew',842,4.9,110,278);
/*!40000 ALTER TABLE `Product` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ShoppingCart`
--

DROP TABLE IF EXISTS `ShoppingCart`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ShoppingCart` (
  `UserID` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `ProductID` int NOT NULL,
  `Qty` int NOT NULL,
  `Modified` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`UserID`,`ProductID`),
  KEY `UserID` (`UserID`),
  KEY `ProductID` (`ProductID`),
  CONSTRAINT `shoppingcart_ibfk_1` FOREIGN KEY (`UserID`) REFERENCES `User` (`UserID`),
  CONSTRAINT `shoppingcart_ibfk_2` FOREIGN KEY (`ProductID`) REFERENCES `Product` (`ProductID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ShoppingCart`
--

LOCK TABLES `ShoppingCart` WRITE;
/*!40000 ALTER TABLE `ShoppingCart` DISABLE KEYS */;
/*!40000 ALTER TABLE `ShoppingCart` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `User`
--

DROP TABLE IF EXISTS `User`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `User` (
  `UserID` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `Password` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `FirstName` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `LastName` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `Phone` int unsigned NOT NULL,
  `Email` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `Address` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `Seller` tinyint(1) DEFAULT '0',
  `Verified` tinyint(1) NOT NULL DEFAULT '0',
  `Created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`UserID`),
  UNIQUE KEY `Email` (`Email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `User`
--

LOCK TABLES `User` WRITE;
/*!40000 ALTER TABLE `User` DISABLE KEYS */;
INSERT INTO `User` VALUES ('ahmadmuhammad','$2a$04$WcKjLdfbzZ.ipq8XIEe5eOA6Sa/Y8nFU4iZ6CVO6woykxZSCRs3km','Ahamad','Muhammad',63986983,'amuhammad23@gmail.com','18 Sungei Tengah Road, Singapore 698974',1,1,'2021-05-09 05:02:48'),('balamuthu','$2a$04$lmHqFPG9w.L1KDc6MbrNZe7sH1kiI9Uz74lL.RbTtx2t5WNQY32vS','Bala','Muthu',92374639,'bmuthu87@gmail.com','100 Neo Tiew Rd, Singapore 719026',1,1,'2021-05-09 05:02:48'),('jasonlim','$2a$04$ImrcsFqzfaLrZo57.rmj3OWR30dYS8y7dXqkU01wz5.PAwEM/Nu3W','Jason','Lim',96783912,'jasonlim@gmail.com','16 Kallang Place #04-19/20, Singapore 339156',1,1,'2021-05-09 05:02:48'),('leematthew','$2a$04$eMCDHGPUyhvaUuYsHeI7QePJSHorOzeI3Br6MTpc0UWbXaRArvZx6','Matthew','Lee',96548324,'matthewlee@gmail.com','15 Woodlands Loop, Singapore 738322',1,1,'2021-05-09 05:02:48'),('ongryan123','$2a$04$AgVFx4UV/bX2kZR7eRKAGeTQQLQ02Q8wpB3ZfTV6InSSC22IMTt4K','Ryan','Ong',92649843,'ryanong95@gmail.com','60 Jalan Penjara, Singapore 149375',0,1,'2021-05-09 05:02:48'),('tanmayling','$2a$04$Rvsgd2SPTAG5xXsNv.MO4e8c2Ui6BBC8zNdE7q/4ZO/I5Lb8wPBMm','May Ling','Tan',83457233,'tanmayling@gmail.com','220 Neo Tiew Crescent, Singapore 718830',0,1,'2021-05-09 05:02:48');
/*!40000 ALTER TABLE `User` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-05-14 14:14:14
