import java.util.*;

	public class temp
{
		
	public static void main(String[] args) 
	{
		float tempInF;
		float tempInC;
		       
		Scanner in = new Scanner(System.in);

		System.out.printf("Please enter the temperature in Celsius : ");
		
		tempInC = in.nextInt();

		tempInF = 32+ 9 * (tempInC/5);
		System.out.printf("The temperature in Fahrenheit = " + tempInF + "\n");

		System.out.printf("Please enter the temperature in Fahrenheit : ");
		tempInF = in.nextInt();

		tempInC = ((tempInF - 32) * 5) / 9;

		System.out.printf("The temperature in Celsius = " + tempInC + "\n");
		           
	}
}

/* OUTPUT
  -------
  

Please enter the temperature in Celsius : 89
The temperature in Fahrenheit = 192.2
Please enter the temperature in Fahrenheit : 95
The temperature in celsius = 35.0


---------------------------------------------------------

DESIGN
--------
Variables:

(float) : both celsius and fahrenheit in float 

INPUT :

asks user to enter the temp in  celsius and fahrenheit
System.out.printf() for statements for output 

PROCESSING:

(used 5 and 9 instead of 100 and 180 for fast process and easier calculation)
 takes celsius amount takes fahrenheit then subtracts 32 first then multiplies by 5 then divide by 9
 for fahrenheit amount divides by 5 and multiplies by 9 then add 32 to get fahrenheit
 
 Printing :
 
 used System.out.printf() to ask temperatures to user and just turned it into a different temp.

 



	 */
	