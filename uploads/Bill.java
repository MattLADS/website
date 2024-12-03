
import java.util.*;

	public class Bill{
		static final double tv_price= 400.00;
		static final double vcr_price = 220;
		static final double remote_price = 35.20;
		static final double cd_price = 300.00;
		static final double tape_price = 150.00;
		static final double tax = 8.25;
	

	public static void main(String[] args) {
		int num_tv ;
		int num_vcr;
		int num_remote;
		int num_cd ;
		int num_tape;
		
		double total_tvprice;
		double total_vcrprice;
		double total_remoteprice;
		double total_cdprice;
		double total_tapeprice;
		double subtotal;
		double tax_amount;
		double total;
		
		Scanner input = new Scanner(System.in);
		
		System.out.print("How many TV's were sold?");
		num_tv = input.nextInt();
		System.out.print("How many VCR's were sold?");
		num_vcr = input.nextInt();
		System.out.print("How many remote controller's were sold?");
		num_remote = input.nextInt();
		System.out.print("How many CD's were sold?");
		num_cd = input.nextInt();
		System.out.print("How many Tape Recorder's were sold?");
		num_tape = input.nextInt();
		
		total_tvprice = tv_price * num_tv;
		total_vcrprice = vcr_price * num_vcr;
		total_remoteprice = remote_price * num_remote;
		total_cdprice = cd_price * num_cd;
		total_tapeprice = tape_price * num_tape;
		
		subtotal = (total_tvprice + total_vcrprice +total_remoteprice+total_cdprice + total_tapeprice);
		tax_amount = (tax/100) * subtotal;
		total = subtotal + tax_amount;
		
		System.out.printf("%-15s %-20s %-15s %-15s\n", "QTY", "DESCRIPTION", "UNIT PRICE", "TOTAL PRICE");
		
		System.out.printf("%-15d %-20s %-15.2f %-15.2f\n",num_tv, "TV", tv_price, total_tvprice);

		System.out.printf("%-15d %-20s %-15.2f %-15.2f\n", num_vcr, "VCR", vcr_price, total_vcrprice);
		
		System.out.printf("%-15d %-20s %-15.2f %-15.2f\n",num_remote,"REMOTE CONTROLLER", remote_price,total_remoteprice);
		
		System.out.printf("%-15d %-20s %-15.2f %-15.2f\n", num_cd, "CD Player", cd_price, total_cdprice);
		
		System.out.printf("%-15d %-20s %-15.2f %-15.2f\n", num_tape, "TAPE RECORDER", tape_price, total_tapeprice);
		
		System.out.printf("\n %37s %.2f\n", "SUBTOTAL $", subtotal );
		
		System.out.printf("%33s %.2f\n", "TAX $", tax_amount);
		
		System.out.printf("%35s %.2f\n", "TOTAL $", total);
		
		input.close();
	}

}
	

/*-------------------------------------------------------------------------
 OUTPUT set 1/
	  
How many TV's were sold?2
How many VCR's were sold?1
How many remote controller's were sold?4
How many CD's were sold?1
How many Tape Recorder's were sold?2
QTY             DESCRIPTION          UNIT PRICE      TOTAL PRICE    
2               TV                   400.00          800.00         
1               VCR                  220.00          220.00         
4               REMOTE CONTROLLER    35.20           140.80         
1               CD Player            300.00          300.00         
2               TAPE RECORDER        150.00          300.00         

                            SUBTOTAL $ 1760.80
                            TAX $ 145.27
                            TOTAL $ 1906.07
                            
 ------------------------------------------------------------------------    
 OUTPUT set 2/
 
How many TV's were sold?3
How many VCR's were sold?0
How many remote controller's were sold?2
How many CD's were sold?0
How many Tape Recorder's were sold?21
QTY             DESCRIPTION          UNIT PRICE      TOTAL PRICE    
3               TV                   400.00          1200.00        
0               VCR                  220.00          0.00           
2               REMOTE CONTROLLER    35.20           70.40          
0               CD Player            300.00          0.00           
21              TAPE RECORDER        150.00          3150.00        

                            SUBTOTAL $ 4420.40
                            TAX $ 364.68
                            TOTAL $ 4785.08
------------------------------------------------------------------------------  
DESIGN
*  //Variables//
* (int) :variables to track the quantities of each items
* (double) :constants for price of each item per single sale
* (double) : variable for total price after entering numbers of sale and price
* (double) : variable for sub-total, tax amount and final total.
*
*   // INPUT//
*	
*  use scanner to read the quantities and price of each item. 
*  System.out.printf() for printing the questions, table and totals.	
*	item.nextInt() to enter the quantity.
*
*  //PROCESSING//
* multiply price with quantity
* add all the prices together for sub-total
* divide tax% by 100 and multiply with sub-total to get tax amount
* add tax amount into sub-total for final total
*
* // PRINTING//
* 
* 	asks the user how many of each item were sold
* 	after user enter the amounts draw a table with QTY, DESCRIPTION, UNIT PRICE and total price
* 	have everything in row and column order starting each name under the previous name
*	sub- total ,tax and final total under the whole table.
*
*
*
*
*
*
*
*
*
*
*
*
*/
