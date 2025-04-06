package models

import (
	"gorm.io/gorm"
)

const (
	// Address Status
	AddressStatusDisabled int8 = -1 // 禁用
	AddressStatusEnabled  int8 = 10 // 启用
)

// Address 定义与数据库表对应的结构体
type Address struct {
	BaseModel
	AdminID   uint   `gorm:"type:int unsigned not null;index;comment:管理员ID" json:"adminId"`
	CountryID uint   `gorm:"type:varchar(100) not null;comment:国家ID" json:"countryId"`
	City      string `gorm:"type:varchar(100) not null;comment:城市" json:"city"`
	Address   string `gorm:"type:varchar(255) not null;comment:详细地址" json:"address"`
	Status    int8   `gorm:"type:tinyint not null;default:10;index;comment:状态(-1:禁用,10:启用)" json:"status"`
}

//func init() {
//	// Create the table if it doesn't exist
//	if err := db.AutoMigrate(&Address{}); err != nil {
//		panic("Failed to auto migrate Address table: " + err.Error())
//	}
//
//	if err := InitAddressTemplate(db); err != nil {
//		panic("Failed to initialize address template : " + err.Error())
//	}
//}

func InitAddressTemplate(db *gorm.DB) error {
	var count int64
	if err := db.Model(&Address{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		addressTemplate := []Address{
			{AdminID: SuperAdminID, CountryID: 2, City: "Los Angeles", Address: "123 Main St, Los Angeles, CA 90001"},
			{AdminID: SuperAdminID, CountryID: 2, City: "New York", Address: "456 Oak Ave, New York, NY 10001"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Chicago", Address: "789 Maple Rd, Chicago, IL 60601"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Houston", Address: "101 Park Ave, Houston, TX 77001"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Phoenix", Address: "222 Elm St, Phoenix, AZ 85001"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Philadelphia", Address: "345 Central Ave, Philadelphia, PA 19101"},
			{AdminID: SuperAdminID, CountryID: 2, City: "San Antonio", Address: "678 Market St, San Antonio, TX 78201"},
			{AdminID: SuperAdminID, CountryID: 2, City: "San Diego", Address: "901 River Rd, San Diego, CA 92101"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Dallas", Address: "112 Broadway, Dallas, TX 75201"},
			{AdminID: SuperAdminID, CountryID: 2, City: "San Jose", Address: "456 Pine St, San Jose, CA 95101"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Austin", Address: "789 Main St, Austin, TX 78701"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Jacksonville", Address: "101 Oak Ave, Jacksonville, FL 32201"},
			{AdminID: SuperAdminID, CountryID: 2, City: "San Francisco", Address: "222 Central St, San Francisco, CA 94101"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Indianapolis", Address: "345 Market St, Indianapolis, IN 46201"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Columbus", Address: "678 Elm St, Columbus, OH 43201"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Charlotte", Address: "901 Pine Ave, Charlotte, NC 28201"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Fort Worth", Address: "112 Central Ave, Fort Worth, TX 76101"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Detroit", Address: "456 River Rd, Detroit, MI 48201"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Memphis", Address: "789 Broadway, Memphis, TN 38101"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Seattle", Address: "101 Central St, Seattle, WA 98101"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Denver", Address: "222 Market St, Denver, CO 80201"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Boston", Address: "345 Elm St, Boston, MA 02101"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Nashville", Address: "678 Pine Ave, Nashville, TN 37201"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Milwaukee", Address: "901 Central Ave, Milwaukee, WI 53201"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Louisville", Address: "112 River Rd, Louisville, KY 40201"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Baltimore", Address: "456 Market St, Baltimore, MD 21201"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Oklahoma City", Address: "789 Central Ave, Oklahoma City, OK 73101"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Portland", Address: "101 Pine St, Portland, OR 97201"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Tampa", Address: "222 Elm St, Tampa, FL 33601"},
			{AdminID: SuperAdminID, CountryID: 2, City: "New Orleans", Address: "345 River Rd, New Orleans, LA 70101"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Tucson", Address: "678 Central St, Tucson, AZ 85701"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Kansas City", Address: "901 Market St, Kansas City, MO 64101"},
			{AdminID: SuperAdminID, CountryID: 2, City: "St. Louis", Address: "112 Oak Ave, St. Louis, MO 63101"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Fresno", Address: "456 Central Ave, Fresno, CA 93701"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Sacramento", Address: "789 Pine St, Sacramento, CA 95801"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Long Beach", Address: "101 Elm St, Long Beach, CA 90801"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Mesa", Address: "222 River Rd, Mesa, AZ 85201"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Arlington", Address: "345 Central St, Arlington, TX 76001"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Oakland", Address: "678 Market St, Oakland, CA 94601"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Raleigh", Address: "901 Central Ave, Raleigh, NC 27601"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Cincinnati", Address: "112 Pine St, Cincinnati, OH 45201"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Virginia Beach", Address: "456 Elm St, Virginia Beach, VA 23451"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Atlanta", Address: "789 River Rd, Atlanta, GA 30301"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Colorado Springs", Address: "101 Market St, Colorado Springs, CO 80901"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Minneapolis", Address: "222 Central Ave, Minneapolis, MN 55401"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Tulsa", Address: "345 Pine St, Tulsa, OK 74101"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Wichita", Address: "678 Central St, Wichita, KS 67201"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Newark", Address: "901 Market St, Newark, NJ 07101"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Pittsburgh", Address: "112 Oak Ave, Pittsburgh, PA 15201"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Lexington", Address: "456 River Rd, Lexington, KY 40501"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Anchorage", Address: "789 Central Ave, Anchorage, AK 99501"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Albuquerque", Address: "101 Pine St, Albuquerque, NM 87101"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Honolulu", Address: "222 Elm St, Honolulu, HI 96801"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Anaheim", Address: "345 Market St, Anaheim, CA 92801"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Aurora", Address: "678 Central St, Aurora, CO 80001"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Bakersfield", Address: "901 Central Ave, Bakersfield, CA 93301"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Riverside", Address: "112 Pine St, Riverside, CA 92501"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Corpus Christi", Address: "456 Elm St, Corpus Christi, TX 78401"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Stockton", Address: "789 River Rd, Stockton, CA 95201"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Toledo", Address: "101 Market St, Toledo, OH 43601"},
			{AdminID: SuperAdminID, CountryID: 2, City: "St. Petersburg", Address: "222 Central Ave, St. Petersburg, FL 33701"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Cincinnati", Address: "345 Pine St, Cincinnati, OH 45201"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Greensboro", Address: "678 Central St, Greensboro, NC 27401"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Plano", Address: "901 Market St, Plano, TX 75001"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Lincoln", Address: "112 Oak Ave, Lincoln, NE 68501"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Arlington", Address: "456 River Rd, Arlington, VA 22201"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Glendale", Address: "789 Central Ave, Glendale, CA 91201"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Akron", Address: "101 Pine St, Akron, OH 44301"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Chula Vista", Address: "222 Elm St, Chula Vista, CA 91901"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Fort Wayne", Address: "345 Market St, Fort Wayne, IN 46801"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Jersey City", Address: "678 Central St, Jersey City, NJ 07301"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Durham", Address: "901 Central Ave, Durham, NC 27701"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Madison", Address: "112 Pine St, Madison, WI 53701"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Lubbock", Address: "456 Elm St, Lubbock, TX 79401"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Baton Rouge", Address: "789 River Rd, Baton Rouge, LA 70801"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Irvine", Address: "101 Market St, Irvine, CA 92601"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Reno", Address: "222 Central Ave, Reno, NV 89501"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Hialeah", Address: "345 Pine St, Hialeah, FL 33010"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Chesapeake", Address: "678 Central St, Chesapeake, VA 23320"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Scottsdale", Address: "901 Market St, Scottsdale, AZ 85250"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Montgomery", Address: "112 Oak Ave, Montgomery, AL 36101"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Yonkers", Address: "456 River Rd, Yonkers, NY 10701"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Glendale", Address: "789 Central Ave, Glendale, AZ 85301"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Orlando", Address: "101 Pine St, Orlando, FL 32801"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Oakland", Address: "222 Elm St, Oakland, CA 94601"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Winston-Salem", Address: "345 Market St, Winston-Salem, NC 27101"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Chandler", Address: "678 Central St, Chandler, AZ 85224"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Laredo", Address: "901 Central Ave, Laredo, TX 78040"},
			{AdminID: SuperAdminID, CountryID: 2, City: "North Las Vegas", Address: "112 Pine St, North Las Vegas, NV 89030"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Buffalo", Address: "456 Elm St, Buffalo, NY 14201"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Irving", Address: "789 River Rd, Irving, TX 75038"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Rochester", Address: "101 Market St, Rochester, NY 14601"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Fort Lauderdale", Address: "222 Central Ave, Fort Lauderdale, FL 33301"},
			{AdminID: SuperAdminID, CountryID: 2, City: "St. Paul", Address: "345 Pine St, St. Paul, MN 55101"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Norfolk", Address: "678 Central St, Norfolk, VA 23501"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Arlington", Address: "901 Market St, Arlington, TX 76010"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Garland", Address: "112 Oak Ave, Garland, TX 75040"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Richmond", Address: "456 River Rd, Richmond, VA 23219"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Boise", Address: "789 Central Ave, Boise, ID 83701"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Spokane", Address: "101 Pine St, Spokane, WA 99201"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Tacoma", Address: "222 Elm St, Tacoma, WA 98401"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Glendale", Address: "345 Market St, Glendale, CA 91201"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Huntington Beach", Address: "678 Central St, Huntington Beach, CA 92646"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Modesto", Address: "901 Central Ave, Modesto, CA 95350"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Fayetteville", Address: "112 Pine St, Fayetteville, NC 28301"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Moreno Valley", Address: "456 Elm St, Moreno Valley, CA 92553"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Pasadena", Address: "789 River Rd, Pasadena, CA 91101"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Des Moines", Address: "101 Market St, Des Moines, IA 50301"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Grand Rapids", Address: "222 Central Ave, Grand Rapids, MI 49501"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Salt Lake City", Address: "345 Pine St, Salt Lake City, UT 84101"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Huntsville", Address: "678 Central St, Huntsville, AL 35801"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Newport News", Address: "901 Market St, Newport News, VA 23601"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Tempe", Address: "112 Oak Ave, Tempe, AZ 85281"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Lancaster", Address: "456 River Rd, Lancaster, CA 93534"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Fort Collins", Address: "789 Central Ave, Fort Collins, CO 80521"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Tallahassee", Address: "101 Pine St, Tallahassee, FL 32301"},
			{AdminID: SuperAdminID, CountryID: 2, City: "El Cajon", Address: "222 Elm St, El Cajon, CA 92020"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Ontario", Address: "345 Market St, Ontario, CA 91761"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Orange", Address: "678 Central St, Orange, CA 92865"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Brownsville", Address: "901 Central Ave, Brownsville, TX 78520"},
			{AdminID: SuperAdminID, CountryID: 2, City: "Overland Park", Address: "112 Pine St, Overland Park, KS 66201"},
		}

		if err := db.CreateInBatches(addressTemplate, len(addressTemplate)).Error; err != nil {
			return err
		}
	}

	return nil
}
