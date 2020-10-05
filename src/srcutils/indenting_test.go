package srcutils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndentingTabs(t *testing.T) {
	assert := require.New(t)

	input := `
	function something() end
	function something() 
		-- HI
	end
	function something() 
		
	end
	function something() 
	
	end
	function something()

	end
`
	expected := `
function something() end
function something() 
	-- HI
end
function something() 
	
end
function something() 

end
function something()

end
`

	assert.Equal(expected, strings.Join(TrimConsistentIndenting(strings.Split(input, "\n")), "\n"))
}

func TestIndenting3Spaces(t *testing.T) {
	assert := require.New(t)

	input := `
   function something() end
   function something() 
      -- HI
   end
   function something() 
      
   end
   function something() 
   
   end
   function something()

   end
`
	expected := `
function something() end
function something() 
   -- HI
end
function something() 
   
end
function something() 

end
function something()

end
`

	assert.Equal(expected, strings.Join(TrimConsistentIndenting(strings.Split(input, "\n")), "\n"))
}

func TestIndenting2Spaces(t *testing.T) {
	assert := require.New(t)

	input := `
  function something() end
  function something() 
    -- HI
  end
  function something() 
    
  end
  function something() 
  
  end
  function something()

  end
`
	expected := `
function something() end
function something() 
  -- HI
end
function something() 
  
end
function something() 

end
function something()

end
`

	assert.Equal(expected, strings.Join(TrimConsistentIndenting(strings.Split(input, "\n")), "\n"))
}

func TestIndentingNotIndentingMix1(t *testing.T) {
	assert := require.New(t)

	input := `
function something() end
	function something() 
		-- HI
	end
`
	expected := `
function something() end
	function something() 
		-- HI
	end
`

	assert.Equal(expected, strings.Join(TrimConsistentIndenting(strings.Split(input, "\n")), "\n"))
}

func TestIndentingNotIndentingMix2(t *testing.T) {
	assert := require.New(t)

	input := `
	function something() end
    function something() 
        -- HI
    end
`
	expected := `
	function something() end
    function something() 
        -- HI
    end
`

	assert.Equal(expected, strings.Join(TrimConsistentIndenting(strings.Split(input, "\n")), "\n"))
}
